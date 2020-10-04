package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

// PrintObjectTable print an object array to writer
func PrintObjectTable(obj interface{}, writer io.Writer) {
	var header []string
	var rows [][]string
	data, err := json.Marshal(obj)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		var sliceMap = []map[string]interface{}{}
		err = json.Unmarshal(data, &sliceMap)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		if len(sliceMap) > 0 {
			var value []string
			for k, v := range sliceMap[0] {
				header = append(header, k)
				value = append(value, fmt.Sprintf("%v", v))
			}
			rows = append(rows, value)
			sliceMap = sliceMap[1:]
			for _, item := range sliceMap {
				value = []string{}
				for _, headerKey := range header {
					v := item[headerKey]
					value = append(value, fmt.Sprintf("%v", v))
				}
				rows = append(rows, value)
			}
		}
	} else {
		objMap := map[string]interface{}{}
		err = json.Unmarshal(data, &objMap)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		var value []string

		for k, v := range objMap {
			header = append(header, k)
			value = append(value, fmt.Sprintf("%v", v))
		}
		rows = append(rows, value)
	}

	outWrite := tablewriter.NewWriter(writer)
	outWrite.SetHeader(header)
	outWrite.AppendBulk(rows)
	outWrite.SetRowLine(true)
	outWrite.SetAutoWrapText(false)
	outWrite.SetAutoFormatHeaders(true)
	outWrite.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	outWrite.SetAlignment(tablewriter.ALIGN_LEFT)
	outWrite.SetCenterSeparator("")
	outWrite.SetColumnSeparator("")
	outWrite.SetRowSeparator("")
	outWrite.SetHeaderLine(false)
	outWrite.SetBorder(false)
	outWrite.SetTablePadding("\t") // pad with tabs
	outWrite.SetNoWhiteSpace(true)
	outWrite.Render()
}
