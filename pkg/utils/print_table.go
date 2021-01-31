package utils

import (
	"fmt"
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

//PrintInterfacesToConsole print a slice object to console std out
func PrintInterfacesToConsole(objs interface{}) {
	table := BuildPrintTabel(objs)
	table.Render()
}

//BuildPrintTabel build print tables for objects
func BuildPrintTabel(objs interface{}) *tablewriter.Table {
	printData, printHeader := interfacesToTableString(objs)
	table := tablewriter.NewWriter(os.Stdout)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	if len(printHeader) > 0 {
		table.SetHeader(printHeader)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderLine(false)
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(printData) // Add Bulk Data
	return table
}

func interfacesToTableString(objs interface{}) ([][]string, []string) {
	var result [][]string
	var header []string
	items := reflect.ValueOf(objs)
	if reflect.ValueOf(objs).Kind() == reflect.Slice {
		for i := 0; i < items.Len(); i++ {
			item := items.Index(i)
			if item.Kind() == reflect.Struct {
				v := reflect.Indirect(item)
				var fieldValues []string
				for j := 0; j < v.NumField(); j++ {
					var value = v.Field(j).Interface()
					if reflect.TypeOf(value).Kind() == reflect.Ptr {
						value = v.Field(j).Elem().Interface()
					}
					fieldValues = append(fieldValues, fmt.Sprintf("%v", value))
					if i == 0 {
						header = append(header, v.Type().Field(j).Name)
					}
				}
				result = append(result, fieldValues)
			}
		}
		return result, header
	} else if reflect.ValueOf(objs).Kind() == reflect.Struct {
		result = getStructTableString(objs)
		return result, header
	} else if reflect.ValueOf(objs).Kind() == reflect.Ptr && reflect.ValueOf(objs).Elem().Kind() == reflect.Struct {
		objs = reflect.ValueOf(objs).Elem().Interface()
		result = getStructTableString(objs)
		return result, header
	}
	return [][]string{}, []string{}
}

func getStructTableString(structObj interface{}) [][]string {
	var result [][]string

	v := reflect.Indirect(reflect.ValueOf(structObj))
	for j := 0; j < v.NumField(); j++ {
		var fieldValues []string
		fieldValues = append(fieldValues, fmt.Sprintf("%s:", v.Type().Field(j).Name))
		var value = v.Field(j).Interface()
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			value = v.Field(j).Elem().Interface()
		}
		fieldValues = append(fieldValues, fmt.Sprintf("%v", value))
		result = append(result, fieldValues)
	}
	return result
}
