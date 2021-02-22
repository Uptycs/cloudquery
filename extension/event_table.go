package extension

import (
	"context"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/extension/aws/cloudtrail"
	"github.com/Uptycs/cloudquery/extension/gcp/cloudlog"
	"sync"
	"time"
)

// EventTable interface for eventing tables
type EventTable interface {
	GetName() string
	GetColumns() []table.ColumnDefinition
	GetGenFunction() table.GenerateFunc
	Start(ctx context.Context, wg *sync.WaitGroup, socket string, timeout time.Duration)
}

var (
	once           sync.Once
	eventTableList []EventTable = make([]EventTable, 0)
)

// GetEventTables return the list of all eventing tables
func GetEventTables() []EventTable {
	once.Do(func() {
		eventTableList = []EventTable{
			&cloudtrail.CloudTrailEventTable{},
			&cloudlog.CloudLogEventTable{},
		}
	})
	return eventTableList
}
