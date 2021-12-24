package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/conf"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"os"
	"reflect"
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
	GetValue(name FieldInfo) (string, error)
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
	var err error
	switch outputFormat {
	case "plain":
		plainOutput(fields, items)
	case "simple":
		err = tableOutput(fields, headerFunction, items, tableStyleSimple, false)
	case "color":
		err = tableOutput(fields, headerFunction, items, tableStyleColor, true)
	default:
		err = tableOutput(fields, headerFunction, items, tableStyleNormal, true)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func plainOutput(fields []FieldInfo, items []Item) {
	for _, item := range items {
		for _, field := range fields {
			value, err := item.GetValue(field)
			if err != nil {
				fmt.Printf("<error: %v>", err)
			} else {
				fmt.Printf("%s ", value)
			}
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

func tableOutput(fields []FieldInfo, headerFunction HeaderFunction, items []Item, style *table.Style, showHeader bool) error {
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
			if value, err := item.GetValue(field); err != nil {
				return err
			} else {
				values[i] = value
			}
		}
		t.AppendRow(values)
	}

	t.Render()

	return nil
}

var stringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

func ExtractFromExpression(expression string, item interface{}) (string, error) {
	rawValue, err := evaluateExpression(expression, item)
	if err != nil {
		return "", err
	}

	reflectedValue := reflect.ValueOf(rawValue)
	if reflectedValue.Kind() == reflect.Ptr {
		if reflectedValue.IsNil() {
			return "", nil
		}
		reflectedValue = reflectedValue.Elem()
		rawValue = reflectedValue.Interface()
	}

	if s, ok := rawValue.(string); ok {
		return fmt.Sprintf("%s", s), nil
	} else if sg, ok := rawValue.(fmt.Stringer); ok {
		return fmt.Sprintf("%s", sg), nil
	}

	if s, ok := stringifyBuiltinType(reflectedValue); ok {
		return s, nil
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if encoder.Encode(rawValue) == nil {
		buf.Truncate(buf.Len() - 1) // Discard the stupid newline
		return buf.String(), nil
	}

	return "", nil
}

func stringifyBuiltinType(reflectedValue reflect.Value) (string, bool) {
	switch reflectedValue.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Uintptr:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.String: // still possible because of type aliases
		return fmt.Sprintf("%v", reflectedValue), true
	}
	return "", false
}

func evaluateExpression(expression string, item interface{}) (interface{}, error) {
	program, err := expr.Compile(expression, func(c *conf.Config) {
		c.Strict = false
	})

	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, item)
	if err != nil {
		// Ignore errors, probably an index out of bounds or something
		return "", nil
	}
	return output, nil
}
