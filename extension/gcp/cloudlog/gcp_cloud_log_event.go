/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package cloudlog

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"

	"context"
	"encoding/json"
	"fmt"
	osquery "github.com/Uptycs/basequery-go"
	extgcp "github.com/Uptycs/cloudquery/extension/gcp"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"sort"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/Uptycs/cloudquery/utilities"

	"cloud.google.com/go/storage"
	"github.com/Uptycs/basequery-go/plugin/table"
	"google.golang.org/api/logging/v2"
)

type ObjectMarker struct {
	modifiedTime time.Time
	key          string
	dirPath      string
}

// CloudLogEventTable implements EventTable interface
type CloudLogEventTable struct {
	// Marker will always be atleast markerDelayMinutes prior to current time
	markerDelayMinutes int
	// Map of bucketName+logName => ObjectMarker
	markerMap map[string]*ObjectMarker
	// objects which we have processed in last 1 hour
	objectCache *cache.Cache
	client      *osquery.ExtensionManagerClient
	ctx         context.Context
}

var (
	MARKER_DELAY_MINUTES  = 120         // 2 Hours
	LOOKBACK_MINUTES      = 120         // 2 Hours
	CACHE_TIMEOUT_MINUTES = 2 * 24 * 60 // 2 Days
	LOOP_TIMER_SECONDS    = 15 * 60     // 15 Minutes
	TABLE_NAME            = "gcp_cloud_log_events"
)

func (cl *CloudLogEventTable) GetName() string {
	return TABLE_NAME
}

// GetColumns returns the list of columns in the table
func (cl *CloudLogEventTable) GetColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("http_request"),
		table.TextColumn("insert_id"),
		table.TextColumn("labels"),
		table.TextColumn("log_name"),
		table.TextColumn("metadata"),
		table.TextColumn("operation"),
		table.TextColumn("proto_payload"),
		table.TextColumn("receive_timestamp"),
		table.TextColumn("resource"),
		table.TextColumn("severity"),
		table.TextColumn("source_location"),
		table.TextColumn("span_id"),
		table.TextColumn("text_payload"),
		table.TextColumn("timestamp"),
		table.TextColumn("trace"),
		table.TextColumn("trace_sampled"),
	}
}

// GetGenFunction return the function which generates data. For event table this function is no-op
func (cl *CloudLogEventTable) GetGenFunction() table.GenerateFunc {
	return cl.CloudLogGenerate
}

func (cl *CloudLogEventTable) initialize(ctx context.Context, socket string, timeout time.Duration) {
	cl.ctx = ctx
	cl.markerDelayMinutes = MARKER_DELAY_MINUTES
	cl.objectCache = cache.New(time.Duration(CACHE_TIMEOUT_MINUTES)*time.Minute, time.Duration(CACHE_TIMEOUT_MINUTES)*time.Minute)
	cl.markerMap = make(map[string]*ObjectMarker)
	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) > 0 {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			for _, bucket := range account.CloudLogStorageBuckets {
				for _, logName := range bucket.LogNames {
					cl.markerMap[bucket.Name+logName] = nil
				}
			}
		}
	}
	cl.client, _ = osquery.NewClient(socket, timeout)
}

// Start run the event loop
func (cl *CloudLogEventTable) Start(ctx context.Context, wg *sync.WaitGroup, socket string, timeout time.Duration) {
	utilities.GetLogger().Info("Starting event loop")
	wg.Add(1)
	defer wg.Done()
	cl.initialize(ctx, socket, timeout)
	timer1 := time.NewTimer(time.Duration(LOOP_TIMER_SECONDS) * time.Second)

	for {
		select {
		case <-ctx.Done():
			// Shutdown
			timer1.Stop()
			return
		case <-timer1.C:
			cl.runEventLoop()
			timer1 = time.NewTimer(time.Duration(LOOP_TIMER_SECONDS) * time.Second)
		}
	}
}

// CloudLogGenerate returns empty row
func (cl *CloudLogEventTable) CloudLogGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	return nil, nil
}

func (cl *CloudLogEventTable) runEventLoop() {
	if len(utilities.ExtConfiguration.ExtConfGcp.Accounts) > 0 {
		for _, account := range utilities.ExtConfiguration.ExtConfGcp.Accounts {
			if !extgcp.ShouldProcessProject(TABLE_NAME, account.ProjectID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": TABLE_NAME,
				"projectID": account.ProjectID,
			}).Info("processing account")
			cl.processAccountLookupEvents(&account)
		}
	}
}

/*
Dir: <logName>/YYYY/MM/DD/
FileName: HH:00:00_HH:59:59_S0.json
*/
func (cl *CloudLogEventTable) getDirPath(logName string, startTime time.Time) string {
	dirStr := fmt.Sprintf("%04d", startTime.Year()) + "/" + fmt.Sprintf("%02d", startTime.Month()) + "/" + fmt.Sprintf("%02d", startTime.Day())
	return logName + "/" + dirStr
}

func getJSONStr(prop interface{}) string {
	bytes, err := json.Marshal(prop)
	if err != nil {
		utilities.GetLogger().Info("unmarshal failed", err.Error())
		return ""
	}
	return string(bytes)
}

func logEntryToEventRow(entry logging.LogEntry) map[string]string {
	event := make(map[string]string)
	event["http_request"] = getJSONStr(entry.HttpRequest)
	event["insert_id"] = entry.InsertId
	event["labels"] = getJSONStr(entry.Labels)
	event["log_name"] = entry.LogName
	event["metadata"] = getJSONStr(entry.Metadata)
	event["operation"] = getJSONStr(entry.Operation)
	event["proto_payload"] = getJSONStr(entry.ProtoPayload)
	event["receive_timestamp"] = entry.ReceiveTimestamp
	event["resource"] = getJSONStr(entry.Resource)
	event["severity"] = entry.Severity
	event["source_location"] = getJSONStr(entry.SourceLocation)
	event["span_id"] = entry.SpanId
	event["text_payload"] = entry.TextPayload
	event["timestamp"] = entry.Timestamp
	event["trace"] = entry.Trace
	event["trace_sampled"] = utilities.GetStringValue(entry.TraceSampled)
	return event
}

func (cl *CloudLogEventTable) processRecords(account *utilities.ExtensionConfigurationGcpAccount, bucket utilities.CloudLogStorageBucket,
	logName string, key string, jsonData string, outEvents []map[string]string) []map[string]string {
	jsonObj := logging.LogEntry{}
	err := json.Unmarshal([]byte(jsonData), &jsonObj)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"projectID": account.ProjectID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"logName":   logName,
			"key":       key,
			"errString": err.Error(),
		}).Error("failed to parse object data")
		return outEvents
	}
	event := logEntryToEventRow(jsonObj)
	if !extgcp.ShouldProcessEvent(TABLE_NAME, account.ProjectID, bucket.Region, event) {
		return outEvents
	}
	return append(outEvents, event)
}

func (cl *CloudLogEventTable) getObjectReader(account *utilities.ExtensionConfigurationGcpAccount, bucket utilities.CloudLogStorageBucket,
	logName string, obj *storage.ObjectAttrs, rc *storage.Reader) (io.Reader, error) {
	objectBytes, err := ioutil.ReadAll(rc)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"projectID": account.ProjectID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"logName":   logName,
			"key":       obj.Name,
			"errString": err.Error(),
		}).Error("failed to read object data")
		return nil, err
	}
	var stringData string
	stringData = string(objectBytes[:])
	if strings.HasSuffix(obj.Name, "gz") {
		reader, err := gzip.NewReader(strings.NewReader(stringData))
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": TABLE_NAME,
				"projectID": account.ProjectID,
				"region":    bucket.Region,
				"task":      "LookupEvents",
				"bucket":    bucket.Name,
				"logName":   logName,
				"key":       obj.Name,
				"errString": err.Error(),
			}).Error("failed to create gzip reader")
			return nil, err
		}
		return reader, nil
	} else {
		return strings.NewReader(stringData), nil
	}
}

func (cl *CloudLogEventTable) processSingleObject(client *storage.Client, account *utilities.ExtensionConfigurationGcpAccount, bucket utilities.CloudLogStorageBucket, logName string, obj *storage.ObjectAttrs) error {
	_, found := cl.objectCache.Get(bucket.Name + obj.Name)
	if found {
		// we have already processed this file
		return nil
	}
	rc, err := client.Bucket(bucket.Name).Object(obj.Name).NewReader(cl.ctx)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"projectID": account.ProjectID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"logName":   logName,
			"key":       obj.Name,
			"errString": err.Error(),
		}).Error("failed to process object")
		return err
	}
	defer rc.Close()

	reader, err := cl.getObjectReader(account, bucket, logName, obj, rc)
	if err != nil {
		return err
	}
	r := bufio.NewReaderSize(reader, 1024*1024)
	line, isPrefix, err := r.ReadLine()
	events := make([]map[string]string, 0)
	for err == nil && !isPrefix {
		lineStr := string(line)
		// Collect events
		events = cl.processRecords(account, bucket, logName, obj.Name, lineStr, events)
		line, isPrefix, err = r.ReadLine()
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": TABLE_NAME,
		"projectID": account.ProjectID,
		"region":    bucket.Region,
		"task":      "LookupEvents",
		"bucket":    bucket.Name,
		"logName":   logName,
		"key":       obj.Name,
	}).Debug("Added events ", len(events))
	// Send events
	cl.client.StreamEvents(TABLE_NAME, events)

	if isPrefix {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"projectID": account.ProjectID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"logName":   logName,
			"key":       obj.Name,
			"errString": "buffer size too small",
		}).Error("failed to read object data")
		return fmt.Errorf("buffer size too small")
	}
	if err != nil && err != io.EOF {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"projectID": account.ProjectID,
			"region":    bucket.Region,
			"task":      "LookupEvents",
			"bucket":    bucket.Name,
			"logName":   logName,
			"key":       obj.Name,
			"errString": err.Error(),
		}).Error("failed to read object data")
		return err
	}
	utilities.GetLogger().Info("Processed file ", bucket.Name+obj.Name)
	cl.objectCache.Add(bucket.Name+obj.Name, true, 0)
	return nil
}

func (cl *CloudLogEventTable) processObjects(client *storage.Client, account *utilities.ExtensionConfigurationGcpAccount,
	bucket utilities.CloudLogStorageBucket, objs []*storage.ObjectAttrs, dirPath string, logName string) {
	currentTime := time.Now()
	currentMarker := cl.markerMap[bucket.Name+logName]
	if currentMarker != nil && currentMarker.dirPath != dirPath {
		// this marker is for different day
		currentMarker = nil
	}
	// Sort objs in ascending order using Updated time
	sort.Slice(objs, func(p, q int) bool {
		return objs[p].Updated.Before(objs[q].Updated)
	})
	for _, obj := range objs {
		if currentMarker == nil && obj.Updated.Before(currentTime.Add(-time.Duration(LOOKBACK_MINUTES)*time.Minute)) {
			// we dont have a marker set, and current file is not within latest 1 hour. Ignore
			utilities.GetLogger().Info("Ignoring file:", bucket.Name+obj.Name)
			continue
		}
		// Process object
		cl.processSingleObject(client, account, bucket, logName, obj)
		// if object is not within latest cl.markerDelayMinutes
		// and if it is modified after current marker, update the marker
		if currentTime.Sub(obj.Updated) >= time.Duration(cl.markerDelayMinutes)*time.Minute {
			if currentMarker == nil || currentMarker.modifiedTime.Before(obj.Updated) {
				// update marker
				newMarker := ObjectMarker{
					modifiedTime: obj.Updated,
					key:          obj.Name,
					dirPath:      dirPath,
				}
				currentMarker = &newMarker
			}
		}
	}
	if currentMarker != nil {
		cl.markerMap[bucket.Name+logName] = currentMarker
	}
}

func (cl *CloudLogEventTable) getStorageServiceForAccount(account *utilities.ExtensionConfigurationGcpAccount) (*storage.Client, string) {
	var projectID string
	var client *storage.Client
	var err error
	if account != nil {
		projectID = account.ProjectID
		client, err = storage.NewClient(cl.ctx, option.WithCredentialsFile(account.KeyFile))
	} else {
		projectID = utilities.DefaultGcpProjectID
		client, err = storage.NewClient(cl.ctx)
	}
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
			"projectId": projectID,
			"errString": err.Error(),
		}).Error("failed to create client")
		return nil, ""
	}
	return client, projectID
}

func (cl *CloudLogEventTable) getObjectList(client *storage.Client, bucketName, dirPath string) []*storage.ObjectAttrs {
	q := storage.Query{
		Prefix: dirPath,
	}
	objList := make([]*storage.ObjectAttrs, 0)
	it := client.Bucket(bucketName).Objects(cl.ctx, &q)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":  TABLE_NAME,
				"bucketName": bucketName,
				"errString":  err.Error(),
			}).Error("failed to iterate objects")
		}
		objList = append(objList, attrs)
	}

	return objList
}

func (cl *CloudLogEventTable) processBucket(account *utilities.ExtensionConfigurationGcpAccount, bucket utilities.CloudLogStorageBucket) {
	utilities.GetLogger().Info("Processing bucket ", account.ProjectID, ":", bucket.Name)
	client, _ := cl.getStorageServiceForAccount(account)
	if client == nil {
		return
	}
	defer client.Close()

	for _, logName := range bucket.LogNames {
		currentTime := time.Now()
		dirPath := cl.getDirPath(logName, currentTime)
		pastDirPath := cl.getDirPath(logName, currentTime.Add(-time.Duration(cl.markerDelayMinutes)*time.Minute))
		if dirPath != pastDirPath {
			// we just moved to new day, but we need to process last few files in past day as well
			storageObjects := cl.getObjectList(client, bucket.Name, pastDirPath)
			cl.processObjects(client, account, bucket, storageObjects, pastDirPath, logName)
		}
		// process current day
		storageObjects := cl.getObjectList(client, bucket.Name, dirPath)
		cl.processObjects(client, account, bucket, storageObjects, dirPath, logName)
	}
	return
}

func (cl *CloudLogEventTable) processAccountLookupEvents(account *utilities.ExtensionConfigurationGcpAccount) {
	if account == nil || len(account.CloudLogStorageBuckets) == 0 {
		return
	}
	_, ok := utilities.TableConfigurationMap[TABLE_NAME]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": TABLE_NAME,
		}).Error("failed to get table configuration")
		return
	}
	for _, bucket := range account.CloudLogStorageBuckets {
		cl.processBucket(account, bucket)
	}

	return
}
