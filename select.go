package sqlbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Select
// Select sql
type SelectBuilder struct {
	where   where
	isBuilt bool
	sql     string
	table   string
	fields  []string
	data    []interface{}
}

func (builder *SelectBuilder) Build() (string, []interface{}) {
	if !builder.isBuilt {
		sql := fmt.Sprintf("select %s from %s", strings.Join(builder.fields, ", "), builder.table)
		if builder.where != nil {
			where, whereData := builder.where.Build()
			sql = fmt.Sprintf("%s where %s", sql, where)
			builder.data = whereData
		}
		builder.sql = sql
		builder.isBuilt = true
	}
	return builder.sql, builder.data
}

func Select(fields ...string) *SelectBuilder {
	builder := &SelectBuilder{
		fields: fields,
	}
	return builder
}

func (builder *SelectBuilder) From(table string) *SelectBuilder {
	builder.table = table
	return builder
}

func (builder *SelectBuilder) getWhere() where {
	if builder.where == nil {
		builder.where = &WhereBuilder{}
	}
	return builder.where
}

// WhereOperate where column > value
// 可以指定操作符的where语句
func (builder *SelectBuilder) WhereOperate(column, operate string, value interface{}) *SelectBuilder {
	builder.getWhere().WhereOperate(column, operate, value)
	return builder
}

// Where where some_column = value
// 普通的where语句
func (builder *SelectBuilder) Where(column string, value interface{}) *SelectBuilder {
	builder.getWhere().Where(column, value)
	return builder
}

func (builder *SelectBuilder) WhereIn(column string, value interface{}) *SelectBuilder {
	builder.getWhere().WhereIn(column, value)
	return builder
}

func (builder *SelectBuilder) WhereBetween(column string, min, max interface{}) *SelectBuilder {
	builder.getWhere().WhereBetween(column, min, max)
	return builder
}

func (builder *SelectBuilder) WhereLike(column string, value interface{}) *SelectBuilder {
	builder.getWhere().WhereLike(column, value)
	return builder
}

func (builder *SelectBuilder) WhereNull(column string) *SelectBuilder {
	builder.getWhere().WhereNull(column)
	return builder
}

func (builder *SelectBuilder) WhereNotNull(column string) *SelectBuilder {
	builder.getWhere().WhereNotNull(column)
	return builder
}

func (builder *SelectBuilder) WhereFunc(f BuilderFunc) *SelectBuilder {
	builder.getWhere().WhereFunc(f)
	return builder
}

func (builder *SelectBuilder) OrWhereOperate(column, operate string, value interface{}) *SelectBuilder {
	builder.getWhere().OrWhereOperate(column, operate, value)
	return builder
}

func (builder *SelectBuilder) OrWhere(column string, value interface{}) *SelectBuilder {
	builder.getWhere().OrWhere(column, value)
	return builder
}

func (builder *SelectBuilder) OrWhereFunc(f BuilderFunc) *SelectBuilder {
	builder.getWhere().OrWhereFunc(f)
	return builder
}

func (builder *SelectBuilder) String() string {
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
