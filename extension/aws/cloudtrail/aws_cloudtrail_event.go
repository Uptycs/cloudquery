/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package cloudtrail

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	osquery "github.com/Uptycs/basequery-go"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
)

type ObjectMarker struct {
	modifiedTime time.Time
	key          string
	prefix       string
}

// CloudTrailEventTable implements EventTable interface
type CloudTrailEventTable struct {
	// Marker will always be atleast markerDelayMinutes prior to current time
	markerDelayMinutes int
	// Map of bucketName => ObjectMarker
	markerMap map[string]*ObjectMarker
	// objects which we have processed in last 1 hour
	objectCache *cache.Cache
	client      *osquery.ExtensionManagerClient
	ctx         context.Context
}

type CloudTrailEventRecords struct {
	Records []map[string]interface{} `json:"Records"`
}

var (
	MARKER_DELAY_MINUTES  = 20
	LOOKBACK_MINUTES      = 20
	CACHE_TIMEOUT_MINUTES = 120
	LOOP_TIMER_SECONDS    = 120
	TABLE_NAME            = "aws_cloudtrail_events"
)

func (ct *CloudTrailEventTable) GetName() string {
	return TABLE_NAME
}

// GetColumns returns the list of columns in the table
func (ct *CloudTrailEventTable) GetColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_id"),
		table.TextColumn("region_code"),
		table.TextColumn("addendum"),
		table.TextColumn("additional_event_data"),
		table.TextColumn("api_version"),
		table.TextColumn("error_code"),
		table.TextColumn("error_message"),
		table.TextColumn("event_category"),
		table.TextColumn("event_id"),
		table.TextColumn("event_name"),
		table.TextColumn("event_source"),
		table.TextColumn("event_time"),
		table.TextColumn("event_version"),
		table.TextColumn("insight_details"),
		table.TextColumn("management_event"),
		table.TextColumn("read_only"),
		table.TextColumn("recipient_account_id"),
		table.TextColumn("request_id"),
		table.TextColumn("request_parameters"),
		table.TextColumn("resources"),
		table.TextColumn("response_elements"),
		table.TextColumn("service_event_details"),
		table.TextColumn("session_credential_from_console"),
		table.TextColumn("edge_device_details"),
		table.TextColumn("shared_event_id"),
		table.TextColumn("source_ip_address"),
		table.TextColumn("tls_details"),
		table.TextColumn("user_agent"),
		table.TextColumn("user_identity"),
		table.TextColumn("vpc_endpoint_id"),
	}
}

// GetGenFunction return the function which generates data. For event table this function is no-op
func (ct *CloudTrailEventTable) GetGenFunction() table.GenerateFunc {
	return ct.LookupEventsGenerate
}

func (ct *CloudTrailEventTable) initialize(ctx context.Context, socket string, timeout time.Duration) {
	ct.ctx = ctx
	ct.markerDelayMinutes = MARKER_DELAY_MINUTES
	ct.objectCache = cache.New(time.Duration(CACHE_TIMEOUT_MINUTES)*time.Minute, time.Duration(CACHE_TIMEOUT_MINUTES)*time.Minute)
	ct.markerMap = make(map[string]*ObjectMarker)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) > 0 {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			for _, bucket := range account.CtS3Buckets {
				ct.markerMap[bucket.Name] = nil
			}
		}
	}
	ct.client, _ = osquery.NewClient(socket, timeout)
}

// Start run the event loop
func (ct *CloudTrailEventTable) Start(ctx context.Context, wg *sync.WaitGroup, socket string, timeout time.Duration) {
	utilities.GetLogger().Info("Starting event loop")
	wg.Add(1)
	defer wg.Done()
	ct.initialize(ctx, socket, timeout)
	timer1 := time.NewTimer(time.Duration(LOOP_TIMER_SECONDS) * time.Second)

	for {
		select {

		case <-ctx.Done():
			// Shutdown
			timer1.Stop()
			return
		case <-timer1.C:
			ct.runEventLoop()
			timer1 = time.NewTimer(time.Duration(LOOP_TIMER_SECONDS) * time.Second)
		}
	}
	return
}

// DescribeInstancesGenerate returns the rows in the table for all configured accounts
func (ct *CloudTrailEventTable) LookupEventsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	return nil, nil
}

func (ct *CloudTrailEventTable) runEventLoop() {
	utilities.GetLogger().Info("Collecting events")
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) > 0 {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_acm_certificate", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": TABLE_NAME,
				"account":   account.ID,
			}).Info("processing account")
			ct.processAccountLookupEvents(&account)
		}
	}
}

func (ct *CloudTrailEventTable) getPrefix(account *utilities.ExtensionConfigurationAwsAccount, bucket utilities.CtS3Bucket, startTime time.Time) string {
	// currentTime := time.Now()
	// pastHour := currentTime.Add(-time.Duration(1 * time.Hour))
	return bucket.Prefix + "/" + bucket.Region + "/" + fmt.Sprintf("%04d", startTime.Year()) + "/" + fmt.Sprintf("%02d", startTime.Month()) + "/" + fmt.Sprintf("%02d", startTime.Day())
}

func (ct *CloudTrailEventTable) processRecords(account *utilities.ExtensionConfigurationAwsAccount, tableConfig *utilities.TableConfig, bucket utilities.CtS3Bucket, key string, jsonData string) error {
	jsonObj := CloudTrailEventRecords{}
	err := json.Unmarshal([]byte(jsonData), &jsonObj)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"account":   account.ID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"prefix":    bucket.Prefix,
			"key":       key,
			"errString": err.Error(),
		}).Error("failed to parse S3 object data")
		return err
	}
	records := jsonObj.Records
	events := make([]map[string]string, 0)
	for _, record := range records {
		event := make(map[string]string)
		for key, value := range record {
			event[utilities.GetSnakeCase(key)] = utilities.GetStringValue(value)
		}
		if !extaws.ShouldProcessEvent(TABLE_NAME, account.ID, bucket.Region, event) {
			continue
		}
		events = append(events, event)
	}
	utilities.GetLogger().WithFields(log.Fields{
		"tableName": TABLE_NAME,
		"account":   account.ID,
		"region":    bucket.Region,
		"task":      "LookupEvents",
		"bucket":    bucket.Name,
		"prefix":    bucket.Prefix,
		"key":       key,
	}).Debug("Added events ", len(events))
	ct.client.StreamEvents(TABLE_NAME, events)
	return nil
}

func (ct *CloudTrailEventTable) getObjectReader(account *utilities.ExtensionConfigurationAwsAccount, bucket utilities.CtS3Bucket, obj types.Object, output *s3.GetObjectOutput) (io.Reader, error) {
	s3objectBytes, err := ioutil.ReadAll(output.Body)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"account":   account.ID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"prefix":    bucket.Prefix,
			"key":       obj.Key,
			"errString": err.Error(),
		}).Error("failed to read S3 object data")
		return nil, err
	}
	var stringData string
	stringData = string(s3objectBytes[:])
	if strings.HasSuffix(*obj.Key, "gz") {
		reader, err := gzip.NewReader(strings.NewReader(stringData))
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": TABLE_NAME,
				"account":   account.ID,
				"region":    bucket.Region,
				"task":      "LookupEvents",
				"bucket":    bucket.Name,
				"prefix":    bucket.Prefix,
				"key":       obj.Key,
				"errString": err.Error(),
			}).Error("failed to create gzip reader")
			return nil, err
		}
		return reader, nil
	} else {
		return strings.NewReader(stringData), nil
	}
}

func (ct *CloudTrailEventTable) processSingleObject(svc *s3.Client, account *utilities.ExtensionConfigurationAwsAccount, tableConfig *utilities.TableConfig, bucket utilities.CtS3Bucket, obj types.Object) error {
	_, found := ct.objectCache.Get(bucket.Name + *obj.Key)
	if found {
		// we have already processed this file
		return nil
	}
	params := s3.GetObjectInput{
		Bucket: &bucket.Name,
		Key:    obj.Key,
	}
	output, err := svc.GetObject(ct.ctx, &params)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"account":   account.ID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"prefix":    bucket.Prefix,
			"key":       obj.Key,
			"errString": err.Error(),
		}).Error("failed to process S3 object")
		return err
	}
	reader, err := ct.getObjectReader(account, bucket, obj, output)
	if err != nil {
		return err
	}
	r := bufio.NewReaderSize(reader, 1024*1024)
	line, isPrefix, err := r.ReadLine()
	for err == nil && !isPrefix {
		lineStr := string(line)
		ct.processRecords(account, tableConfig, bucket, *obj.Key, lineStr)
		line, isPrefix, err = r.ReadLine()
	}
	if isPrefix {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"account":   account.ID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"prefix":    bucket.Prefix,
			"key":       *obj.Key,
			"errString": "buffer size too small",
		}).Error("failed to read S3 object data")
		return fmt.Errorf("buffer size too small")
	}
	if err != nil && err != io.EOF {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"account":   account.ID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"prefix":    bucket.Prefix,
			"key":       *obj.Key,
			"errString": err.Error(),
		}).Error("failed to read S3 object data")
		return err
	}
	utilities.GetLogger().Info("Processed file ", bucket.Name+*obj.Key)
	ct.objectCache.Add(bucket.Name+*obj.Key, true, 0)
	return nil
}

func (ct *CloudTrailEventTable) processObjects(svc *s3.Client, account *utilities.ExtensionConfigurationAwsAccount, tableConfig *utilities.TableConfig, bucket utilities.CtS3Bucket, objs []types.Object, prefix string) {
	currentTime := time.Now()
	currentMarker := ct.markerMap[bucket.Name]
	if currentMarker != nil && currentMarker.prefix != prefix {
		// this marker is for different prefix
		currentMarker = nil
	}
	// Sort objs in ascending order using LastModified time
	sort.Slice(objs, func(p, q int) bool {
		return objs[p].LastModified.Before(*objs[q].LastModified)
	})
	for _, obj := range objs {
		if currentMarker == nil && obj.LastModified.Before(currentTime.Add(-time.Duration(time.Duration(LOOKBACK_MINUTES)*time.Minute))) {
			// we dont have a marker set, and current file is not within latest 1 hour. Ignore
			continue
		}
		// Process object
		ct.processSingleObject(svc, account, tableConfig, bucket, obj)
		// if object is not within latest ct.markerDelayMinutes
		// and if it is modified after current marker, update the marker
		if currentTime.Sub(*obj.LastModified) >= time.Duration(time.Duration(ct.markerDelayMinutes)*time.Minute) {
			if currentMarker == nil || currentMarker.modifiedTime.Before(*obj.LastModified) {
				// update marker
				newMarker := ObjectMarker{
					modifiedTime: *obj.LastModified,
					key:          *obj.Key,
					prefix:       prefix,
				}
				currentMarker = &newMarker
			}
		}
	}
	if currentMarker != nil {
		ct.markerMap[bucket.Name] = currentMarker
	}
}

func (ct *CloudTrailEventTable) getS3Objects(svc *s3.Client, accountId string, bucket utilities.CtS3Bucket, prefix string) []types.Object {
	s3Objects := make([]types.Object, 0)
	var startAfter *string = nil
	if ct.markerMap[bucket.Name] != nil {
		startAfter = &ct.markerMap[bucket.Name].key
	}
	params := s3.ListObjectsV2Input{
		Bucket:            &bucket.Name,
		Prefix:            &prefix,
		ContinuationToken: nil,
		StartAfter:        startAfter,
	}
	paginator := s3.NewListObjectsV2Paginator(svc, &params)

	for {
		page, err := paginator.NextPage(ct.ctx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": TABLE_NAME,
				"account":   accountId,
				"region":    bucket.Region,
				"task":      "LookupEvents",
				"errString": err.Error(),
			}).Error("failed to process region")
			return s3Objects
		}
		for _, obj := range page.Contents {
			s3Objects = append(s3Objects, obj)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return s3Objects
}

func (ct *CloudTrailEventTable) processBucket(account *utilities.ExtensionConfigurationAwsAccount, tableConfig *utilities.TableConfig, bucket utilities.CtS3Bucket) {
	utilities.GetLogger().Info("Processing bucket ", account.ID, ":", bucket.Name)
	sess, err := extaws.GetAwsConfig(account, bucket.Region)
	if err != nil {
		return
	}
	accountId := account.ID
	svc := s3.NewFromConfig(*sess)
	currentTime := time.Now()
	prefix := ct.getPrefix(account, bucket, currentTime)
	pastPrefix := ct.getPrefix(account, bucket, currentTime.Add(-time.Duration(time.Duration(ct.markerDelayMinutes)*time.Minute)))
	if prefix != pastPrefix {
		// we just moved to new day, but we need to process last few files in past day as well
		s3Objects := make([]types.Object, 0)
		results := ct.getS3Objects(svc, accountId, bucket, pastPrefix)
		for _, obj := range results {
			s3Objects = append(s3Objects, obj)
		}
		ct.processObjects(svc, account, tableConfig, bucket, s3Objects, pastPrefix)
	}
	// process current day
	s3Objects := make([]types.Object, 0)
	results := ct.getS3Objects(svc, accountId, bucket, prefix)
	for _, obj := range results {
		s3Objects = append(s3Objects, obj)
	}
	ct.processObjects(svc, account, tableConfig, bucket, s3Objects, prefix)
	return
}

func (ct *CloudTrailEventTable) processAccountLookupEvents(account *utilities.ExtensionConfigurationAwsAccount) {
	if account == nil || len(account.CtS3Buckets) == 0 {
		return
	}
	tableConfig, ok := utilities.TableConfigurationMap[TABLE_NAME]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
		}).Error("failed to get table configuration")
		return
	}
	for _, bucket := range account.CtS3Buckets {
		ct.processBucket(account, tableConfig, bucket)
	}

	return
}
