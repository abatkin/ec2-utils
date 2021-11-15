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
	if options.OutputFormat == "plain" {
		plainOutput(fields, items)
	} else {
		tableOutput(fields, headerFunction, items)
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

func tableOutput(fields []string, headerFunction HeaderFunction, items []Item) {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.SetOutputMirror(os.Stdout)
	headerFields := make([]interface{}, len(fields))
	for i, field := range fields {
		headerFields[i] = headerFunction(field)
	}
	t.AppendHeader(headerFields)

	for _, item := range items {
		values := make([]interface{}, len(fields))
		for i, field := range fields {
			values[i] = item.GetValue(field)
		}
		t.AppendRow(values)
	}

	t.Render()
}
