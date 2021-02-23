/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package aws

import (
	"context"
	"github.com/Uptycs/basequery-go/plugin/table"
)

// ShouldProcessAccount returns false if given account is not supposed to be processed for given table
// Default implementation is no-op (return true always). Add custom logic here if required
func ShouldProcessAccount(tableName string, accountId string) bool {
	return true
}

// ShouldProcessRegion returns false if given region for given account is not supposed to be processed for given table
// Default implementation is no-op (return true always). Add custom logic here if required
func ShouldProcessRegion(tableName string, accountId string, region string) bool {
	return true
}

// ShouldProcessRow returns false if given row is not supposed to be processed for given table
// Default implementation is no-op (return true always). Add custom logic here if required
func ShouldProcessRow(osqCtx context.Context, queryContext table.QueryContext, tableName string, accountId string, region string, row map[string]interface{}) bool {
	return true
}

// ShouldProcessEvent returns false if given event is not supposed to be processed for given table
// Default implementation is no-op (return true always). Add custom logic here if required
func ShouldProcessEvent(tableName string, accountId string, region string, row map[string]string) bool {
	return true
}
