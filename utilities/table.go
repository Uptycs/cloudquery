package utilities

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Table holds data rows. Each row is map of key->value pair
type Table struct {
	Rows                     []map[string]interface{}
	MaxLevel                 int
	ParsedAttributeConfigMap map[string]ParsedAttributeConfig
}

// NewTable creates a table from given data (in json form) and table configuration
func NewTable(jsonStr []byte, tableConfig *TableConfig) Table {
	tab := Table{}
	tab.init(jsonStr, tableConfig.MaxLevel, tableConfig.getParsedAttributeConfigMap())
	return tab
}

func (tab *Table) init(jsonStr []byte, maxLevel int, parsedAttributeConfigMap map[string]ParsedAttributeConfig) {
	var fields interface{}
	json.Unmarshal(jsonStr, &fields)
	tab.MaxLevel = maxLevel + 1
	tab.ParsedAttributeConfigMap = parsedAttributeConfigMap
	switch fields.(type) {
	case map[string]interface{}:
		tab.flattenMap(0, "", fields.(map[string]interface{}))
	case []interface{}:
		tab.flattenList(0, "", fields.([]interface{}))
	case reflect.Value:
		tab.flattenValue(0, "", fields.(reflect.Value))
	default:
		GetLogger().WithFields(log.Fields{
			"type": reflect.TypeOf(fields),
			"kind": reflect.ValueOf(fields).Kind(),
		}).Warn("Invalid object")
	}

	// fmt.Printf("Flattening fieldMap of size %d\n", len(fieldMap))
	//tab.flattenMap(0, "", fieldMap)
	//tab.print()
}

func (tab *Table) print() {
	for _, row := range tab.Rows {
		GetLogger().Info("===========================")
		for key, value := range row {
			logStr := fmt.Sprintf("%s=%v", key, value)
			GetLogger().Info(logStr)
		}
	}
}

func (tab *Table) addAttribute(name string, value interface{}) {
	// Add attribute only if it is configured
	if attrConfig, ok := tab.ParsedAttributeConfigMap[name]; ok {
		if attrConfig.Enabled {
			if len(tab.Rows) == 0 {
				row := make(map[string]interface{})
				tab.Rows = append(tab.Rows, row)
			}
			for _, item := range tab.Rows {
				item[name] = value
			}
		}
	}
}

func (tab *Table) addRows(newRows []map[string]interface{}) {
	if len(newRows) == 0 {
		// nothing to add
		return
	}

	for _, row := range newRows {
		tab.Rows = append(tab.Rows, row)
	}
}

func (tab *Table) addRowsAndFlatten(newRows []map[string]interface{}) {
	if len(tab.Rows) == 0 {
		tab.Rows = newRows
		return
	} else if len(newRows) == 0 {
		// nothing to flatten
		return
	}
	mergedRows := make([]map[string]interface{}, 0)
	for _, item1 := range tab.Rows {
		for _, item2 := range newRows {
			row := make(map[string]interface{})
			// Add attributes from existing rows
			for key1, value1 := range item1 {
				row[key1] = value1
			}
			// Add attributes from new rows
			for key2, value2 := range item2 {
				row[key2] = value2
			}
			mergedRows = append(mergedRows, row)
		}
	}
	tab.Rows = mergedRows
}

func getKey(prefix, key string) string {
	if len(prefix) != 0 {
		return prefix + "_" + key
	}
	return key
}

// Flatten takes a map and returns a new one where nested maps are replaced
// by dot-delimited keys.
func (tab *Table) flattenMap(level int, prefix string, m map[string]interface{}) {
	for k, v := range m {
		if _, ok := tab.ParsedAttributeConfigMap[getKey(prefix, k)]; ok {
			byteArr, err := json.Marshal(v)
			if err == nil {
				tab.addAttribute(getKey(prefix, k), string(byteArr))
			}
		}
		if tab.MaxLevel > 0 && level >= tab.MaxLevel {
			// Don't flatten further
			continue
		}
		switch child := v.(type) {
		case map[string]interface{}:
			tab.flattenMap(level+1, getKey(prefix, k), child)
		case []interface{}:
			tab.flattenList(level+1, getKey(prefix, k), child)
		case reflect.Value:
			tab.flattenValue(level, getKey(prefix, k), child)
		default:
			tab.addAttribute(getKey(prefix, k), v)
		}
	}
}

func (tab *Table) flattenList(level int, prefix string, list []interface{}) {
	newTable := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
	for _, value := range list {
		if _, ok := tab.ParsedAttributeConfigMap[prefix]; ok {
			scalarTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			byteArr, err := json.Marshal(value)
			if err == nil {
				scalarTab.addAttribute(prefix, string(byteArr))
				newTable.addRows(scalarTab.Rows)
			}
		}
		if tab.MaxLevel > 0 && level >= tab.MaxLevel {
			// Don't flatten further
			continue
		}
		switch child := value.(type) {
		case map[string]interface{}:
			mapTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			mapTab.flattenMap(level+1, prefix, child)
			newTable.addRows(mapTab.Rows)
			//tab.addRowsAndFlatten(newTab.Rows)
		case []interface{}:
			listTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			listTab.flattenList(level+1, prefix, child)
			newTable.addRows(listTab.Rows)
		case reflect.Value:
			valTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			valTab.flattenValue(level, prefix, child)
			newTable.addRows(valTab.Rows)
		default:
			scalarTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			scalarTab.addAttribute(prefix, value)
			newTable.addRows(scalarTab.Rows)
		}
	}
	tab.addRowsAndFlatten(newTable.Rows)
}

func getAdjustedValue(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

func (tab *Table) addAttributeForPrefix(prefix string, value reflect.Value) {
	if _, ok := tab.ParsedAttributeConfigMap[prefix]; ok {
		byteArr, err := json.Marshal(value)
		if err == nil {
			tab.addAttribute(prefix, string(byteArr))
		}
	}
}

func (tab *Table) flattenValue(level int, prefix string, value reflect.Value) {
	value = getAdjustedValue(value)
	tab.addAttributeForPrefix(prefix, value)

	if tab.MaxLevel > 0 && level >= tab.MaxLevel {
		// Don't flatten further
		return
	}

	switch value.Kind() {
	case reflect.Struct:
		var names []string
		for i := 0; i < value.Type().NumField(); i++ {
			name := value.Type().Field(i).Name
			f := value.Field(i)
			if name[0:1] == strings.ToLower(name[0:1]) {
				continue // ignore unexported fields
			}
			if (f.Kind() == reflect.Ptr || f.Kind() == reflect.Slice || f.Kind() == reflect.Map) && f.IsNil() {
				continue // ignore unset fields
			}
			names = append(names, name)
		}
		fieldMap := make(map[string]interface{}, 0)
		for _, n := range names {
			val := value.FieldByName(n)
			fieldMap[n] = val
		}
		tab.flattenMap(level+1, prefix, fieldMap)
	case reflect.Slice:
		fieldList := make([]interface{}, 0)
		for i := 0; i < value.Len(); i++ {
			fieldList = append(fieldList, value.Index(i))
		}
		tab.flattenList(level+1, prefix, fieldList)
	case reflect.Map:
		fieldMap := make(map[string]interface{}, 0)
		for _, k := range value.MapKeys() {
			fieldMap[k.String()] = value.MapIndex(k)
		}
		tab.flattenMap(level+1, prefix, fieldMap)
	default:
		tab.addAttribute(prefix, value.Interface())
	}
}
