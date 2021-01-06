package utilities

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Table struct {
	Rows                     []map[string]interface{}
	MaxLevel                 int
	ParsedAttributeConfigMap map[string]ParsedAttributeConfig
}

func (tab *Table) Init(jsonStr []byte, maxLevel int, parsedAttributeConfigMap map[string]ParsedAttributeConfig) {
	var fieldMap map[string]interface{}
	json.Unmarshal(jsonStr, &fieldMap)
	tab.MaxLevel = maxLevel
	tab.ParsedAttributeConfigMap = parsedAttributeConfigMap
	// fmt.Printf("Flattening fieldMap of size %d\n", len(fieldMap))
	tab.flattenMap(0, "", fieldMap)
	//tab.Print()
}

func (tab *Table) Print() {
	for index, row := range tab.Rows {
		fmt.Printf("[%d] =========================== \n", index)
		for key, value := range row {
			fmt.Printf("%s=%v\n", key, value)
		}
	}
}

func (tab *Table) AddAttribute(name string, value interface{}) {
	// Add attribute only if it is configured
	if _, ok := tab.ParsedAttributeConfigMap[name]; ok {
		if len(tab.Rows) == 0 {
			row := make(map[string]interface{})
			tab.Rows = append(tab.Rows, row)
		}
		for _, item := range tab.Rows {
			item[name] = value
		}
	}
}

func (tab *Table) AddRows(newRows []map[string]interface{}) {
	if len(newRows) == 0 {
		// nothing to add
		return
	}

	for _, row := range newRows {
		tab.Rows = append(tab.Rows, row)
	}
}

func (tab *Table) AddRowsAndFlatten(newRows []map[string]interface{}) {
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
	} else {
		return key
	}
}

// Flatten takes a map and returns a new one where nested maps are replaced
// by dot-delimited keys.
func (tab *Table) flattenMap(level int, prefix string, m map[string]interface{}) {
	for k, v := range m {
		if _, ok := tab.ParsedAttributeConfigMap[getKey(prefix, k)]; ok {
			byteArr, err := json.Marshal(v)
			if err == nil {
				tab.AddAttribute(getKey(prefix, k), string(byteArr))
			}
		}
		if tab.MaxLevel > 0 && level >= tab.MaxLevel {
			// Don't flatten further
			// fmt.Printf("Not Flattening map for field %s. Level:%d, MaxLevel:%d\n", prefix, level, tab.MaxLevel)
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
			tab.AddAttribute(getKey(prefix, k), v)
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
				scalarTab.AddAttribute(prefix, string(byteArr))
				newTable.AddRows(scalarTab.Rows)
			}
		}
		if tab.MaxLevel > 0 && level >= tab.MaxLevel {
			// Don't flatten further
			//fmt.Println("Not Flattening list for field " + prefix)
			continue
		}
		switch child := value.(type) {
		case map[string]interface{}:
			mapTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			mapTab.flattenMap(level+1, prefix, child)
			newTable.AddRows(mapTab.Rows)
			//tab.AddRowsAndFlatten(newTab.Rows)
		case []interface{}:
			listTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			listTab.flattenList(level+1, prefix, child)
			newTable.AddRows(listTab.Rows)
		case reflect.Value:
			valTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			valTab.flattenValue(level, prefix, child)
			newTable.AddRows(valTab.Rows)
		default:
			scalarTab := Table{MaxLevel: tab.MaxLevel, ParsedAttributeConfigMap: tab.ParsedAttributeConfigMap}
			scalarTab.AddAttribute(prefix, value)
			newTable.AddRows(scalarTab.Rows)
		}
	}
	tab.AddRowsAndFlatten(newTable.Rows)
}

func (tab *Table) flattenValue(level int, prefix string, value reflect.Value) {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if _, ok := tab.ParsedAttributeConfigMap[prefix]; ok {
		byteArr, err := json.Marshal(value)
		if err == nil {
			tab.AddAttribute(prefix, string(byteArr))
		}
	}
	if tab.MaxLevel > 0 && level >= tab.MaxLevel {
		// Don't flatten further
		//fmt.Println("Not Flattening value for field " + prefix)
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
		tab.AddAttribute(prefix, value.Interface())
	}
}
