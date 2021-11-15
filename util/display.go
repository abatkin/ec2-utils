package util

import (
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
)
import "github.com/jedib0t/go-pretty/v6/table"

type DisplayOptions struct {
	OutputFormat string
	Fields       []string
}

func BuildDisplayOptions(rootCmd *cobra.Command) *DisplayOptions {
	displayOptions := &DisplayOptions{}

	rootCmd.PersistentFlags().StringVar(&displayOptions.OutputFormat, "output", "table", "Output format")
	rootCmd.PersistentFlags().StringSliceVar(&displayOptions.Fields, "field", []string{}, "List of fields")

	return displayOptions
}

type Field struct {
	FieldName string
	Heading   string
}

type Item interface {
	GetValue(name string) string
}

func (*DisplayOptions) Render(fields []Field, items []Item) {
	t := table.NewWriter()
	t.Style().Format.Header = text.FormatDefault
	t.SetOutputMirror(os.Stdout)
	headerFields := make([]interface{}, len(fields))
	for i, field := range fields {
		headerFields[i] = field.Heading
	}
	t.AppendHeader(headerFields)

	for _, item := range items {
		values := make([]interface{}, len(fields))
		for i, field := range fields {
			values[i] = item.GetValue(field.FieldName)
		}
		t.AppendRow(values)
	}

	t.Render()
}
