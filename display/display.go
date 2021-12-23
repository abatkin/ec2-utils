package display

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
)
import "github.com/jedib0t/go-pretty/v6/table"

type Options struct {
	OutputFormat string
	Fields       []string
}

type Item interface {
	GetValue(name string) string
}

type HeaderFunction func(field string) string

func BuildDisplayOptions(rootCmd *cobra.Command) *Options {
	displayOptions := &Options{}

	rootCmd.PersistentFlags().StringVar(&displayOptions.OutputFormat, "output", "table", "Output format")
	rootCmd.PersistentFlags().StringSliceVar(&displayOptions.Fields, "field", []string{}, "List of fields")

	return displayOptions
}

func (options *Options) Render(fields []string, headerFunction HeaderFunction, items []Item) {
	switch options.OutputFormat {
	case "plain": plainOutput(fields, items)
	case "simple": tableOutput(fields, headerFunction, items, tableStyleSimple, false)
	case "color": tableOutput(fields, headerFunction, items, tableStyleColor, true)
	default: tableOutput(fields, headerFunction, items, tableStyleNormal, true)
	}
}

func plainOutput(fields []string, items []Item) {
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

func tableOutput(fields []string, headerFunction HeaderFunction, items []Item, style *table.Style, showHeader bool) {
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
