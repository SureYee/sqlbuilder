package sqlbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type InsertBuilder struct {
	isBuilt bool
	sql     string
	table   string
	field   []string
	data    []interface{}
}

func Insert(table string) *InsertBuilder {
	return &InsertBuilder{
		table: table,
	}
}

func (builder *InsertBuilder) Fields(fields ...string) *InsertBuilder {
	builder.field = append(builder.field, fields...)
	return builder
}

func (builder *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	builder.data = append(builder.data, values...)
	return builder
}

func (builder *InsertBuilder) Map(mapData map[string]interface{}) *InsertBuilder {
	for column, value := range mapData {
		builder.field = append(builder.field, column)
		builder.data = append(builder.data, value)
	}

	return builder
}

func (builder *InsertBuilder) Build() (string, []interface{}) {
	sql := fmt.Sprintf("insert into %s", builder.table)

	if len(builder.field) > 0 {
		sql = fmt.Sprintf("%s (%s)", sql, strings.Join(builder.field, ", "))
	}

	replace := make([]string, len(builder.data))
	for i := 0; i < len(builder.data); i++ {
		replace[i] = "?"
	}

	sql = fmt.Sprintf("%s values (%s)", sql, strings.Join(replace, ", "))
	builder.sql = sql
	builder.isBuilt = true
	return builder.sql, builder.data
}

func (builder InsertBuilder) String() string {
	sql, data := builder.Build()
	index := 0
	newSql := make([]rune, 0, len(data))
	getData := func(data []interface{}, index int) string {
		if index > (len(data) - 1) {
			return ""
		}

		datum := data[index]

		v := reflect.ValueOf(datum)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.Itoa(int(v.Int()))
		case reflect.String:
			return "\"" + strings.ReplaceAll(strings.ReplaceAll(v.String(), "\\", "\\\\"), "\"", "\\\"") + "\""
		default:
			return ""
		}
	}

	for _, sqlRune := range sql {
		if sqlRune == rune('?') {
			// 将rune替换成 值
			repData := getData(data, index)
			newSql = append(newSql, []rune(repData)...)
			index++
		} else {
			newSql = append(newSql, sqlRune)
		}
	}
	return string(newSql)
}
