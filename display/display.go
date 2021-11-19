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
	case "simple": tableOutput(fields, headerFunction, items, "simple")
	default: tableOutput(fields, headerFunction, items, "table")
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

func tableOutput(fields []string, headerFunction HeaderFunction, items []Item, tableType string) {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.SetOutputMirror(os.Stdout)

	if tableType == "table" {
		headerFields := make([]interface{}, len(fields))
		for i, field := range fields {
			headerFields[i] = headerFunction(field)
		}
		t.AppendHeader(headerFields)
	} else {
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateFooter = false
		t.Style().Options.SeparateHeader = false
		t.Style().Options.SeparateRows = false
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
