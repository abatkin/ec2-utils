package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/lysu/go-el"
	"github.com/spf13/cobra"
	"os"
	"strings"
)
import "github.com/jedib0t/go-pretty/v6/table"

type FieldInfo struct {
	Name       string
	Expression string
}

type Options struct {
	OutputFormat string
	Fields       []string
}

type Item interface {
	GetValue(name FieldInfo) string
}

type HeaderFunction func(field FieldInfo) string

func BuildDisplayOptions(rootCmd *cobra.Command) *Options {
	displayOptions := &Options{}

	rootCmd.PersistentFlags().StringVar(&displayOptions.OutputFormat, "output", "table", "Output format")
	rootCmd.PersistentFlags().StringSliceVar(&displayOptions.Fields, "field", []string{}, "List of fields")

	return displayOptions
}

func ParseFields(rawFields []string) []FieldInfo {
	infos := make([]FieldInfo, len(rawFields))
	for idx, rawField := range rawFields {
		infos[idx] = parseField(rawField)
	}
	return infos
}

func parseField(field string) FieldInfo {
	parts := strings.SplitN(field, "=", 2)
	if len(parts) == 1 {
		return FieldInfo{
			Name:       parts[0],
			Expression: parts[0],
		}
	} else {
		return FieldInfo{
			Name:       parts[0],
			Expression: parts[1],
		}
	}
}

func Render(fields []FieldInfo, outputFormat string, headerFunction HeaderFunction, items []Item) {
	switch outputFormat {
	case "plain":
		plainOutput(fields, items)
	case "simple":
		tableOutput(fields, headerFunction, items, tableStyleSimple, false)
	case "color":
		tableOutput(fields, headerFunction, items, tableStyleColor, true)
	default:
		tableOutput(fields, headerFunction, items, tableStyleNormal, true)
	}
}

func plainOutput(fields []FieldInfo, items []Item) {
	for _, item := range items {
		for _, field := range fields {
			value := item.GetValue(field)
			fmt.Printf("%s ", value)
		}
		fmt.Println()
	}
}

var tableStyleNormal = &table.StyleDefault
var tableStyleColor = &table.StyleColoredBlackOnBlueWhite
var tableStyleSimple = &table.StyleDefault

func init() {
	tableStyleNormal.Format.Header = text.FormatDefault

	tableStyleColor.Format.Header = text.FormatDefault

	tableStyleSimple.Options.DrawBorder = false
	tableStyleSimple.Options.SeparateColumns = false
	tableStyleSimple.Options.SeparateFooter = false
	tableStyleSimple.Options.SeparateHeader = false
	tableStyleSimple.Options.SeparateRows = false
}

func tableOutput(fields []FieldInfo, headerFunction HeaderFunction, items []Item, style *table.Style, showHeader bool) {
	t := table.NewWriter()
	t.SetStyle(*style)
	t.SetOutputMirror(os.Stdout)

	if showHeader {
		headerFields := make([]interface{}, len(fields))
		for i, field := range fields {
			headerFields[i] = headerFunction(field)
		}
		t.AppendHeader(headerFields)
	}

	for _, item := range items {
		values := make([]interface{}, len(fields))
		for i, field := range fields {
			values[i] = item.GetValue(field)
		}
		t.AppendRow(values)
	}

	t.Render()
}

func ExtractFromExpression(expression string, item interface{}) string {
	exp := el.Expression(expression)
	result, err := exp.Execute(item)
	if err != nil {
		return err.Error()
	}

	if result.IsNil() {
		return ""
	}

	switch {
	case result.IsNil():
		return ""
	case result.IsBool():
		fallthrough
	case result.IsNumber():
		fallthrough
	case result.IsString():
		return result.String()
	}

	rawValue := result.Interface()

	if stringer, ok := rawValue.(fmt.Stringer); ok {
		return stringer.String()
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if encoder.Encode(rawValue) == nil {
		buf.Truncate(buf.Len() - 1) // Discard the stupid newline
		return buf.String()
	}

	return ""
}
