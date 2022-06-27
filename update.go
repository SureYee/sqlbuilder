package sqlbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type UpdateBuilder struct {
	isBuilt   bool
	sql       string
	table     string
	where     *WhereBuilder
	fieldData map[string]interface{}
	data      []interface{}
}

func Update(table string) *UpdateBuilder {
	return &UpdateBuilder{
		table:     table,
		fieldData: make(map[string]interface{}),
	}
}

func (builder *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	builder.fieldData[column] = value
	return builder
}

func (builder *UpdateBuilder) Map(mapData map[string]interface{}) *UpdateBuilder {
	for column, value := range mapData {
		builder.Set(column, value)
	}
	return builder
}

func (builder *UpdateBuilder) getWhere() WhereInterface {
	if builder.where == nil {
		builder.where = &WhereBuilder{}
	}
	return builder.where
}

// WhereOperate where column > value
// 可以指定操作符的where语句
func (builder *UpdateBuilder) WhereOperate(column, operate string, value interface{}) *UpdateBuilder {
	builder.getWhere().WhereOperate(column, operate, value)
	return builder
}

// Where where some_column = value
// 普通的where语句
func (builder *UpdateBuilder) Where(column string, value interface{}) *UpdateBuilder {
	builder.getWhere().Where(column, value)
	return builder
}

func (builder *UpdateBuilder) WhereIn(column string, value interface{}) *UpdateBuilder {
	builder.getWhere().WhereIn(column, value)
	return builder
}

func (builder *UpdateBuilder) WhereBetween(column string, min, max interface{}) *UpdateBuilder {
	builder.getWhere().WhereBetween(column, min, max)
	return builder
}

func (builder *UpdateBuilder) WhereLike(column string, value interface{}) *UpdateBuilder {
	builder.getWhere().WhereLike(column, value)
	return builder
}

func (builder *UpdateBuilder) WhereNull(column string) *UpdateBuilder {
	builder.getWhere().WhereNull(column)
	return builder
}

func (builder *UpdateBuilder) WhereNotNull(column string) *UpdateBuilder {
	builder.getWhere().WhereNotNull(column)
	return builder
}

func (builder *UpdateBuilder) WhereFunc(f BuilderFunc) *UpdateBuilder {
	builder.getWhere().WhereFunc(f)
	return builder
}

func (builder *UpdateBuilder) OrWhereOperate(column, operate string, value interface{}) *UpdateBuilder {
	builder.getWhere().OrWhereOperate(column, operate, value)
	return builder
}

func (builder *UpdateBuilder) OrWhere(column string, value interface{}) *UpdateBuilder {
	builder.getWhere().OrWhere(column, value)
	return builder
}

func (builder *UpdateBuilder) OrWhereFunc(f BuilderFunc) *UpdateBuilder {
	builder.getWhere().OrWhereFunc(f)
	return builder
}

func (builder *UpdateBuilder) String() string {
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

func (builder *UpdateBuilder) Build() (string, []interface{}) {
	if !builder.isBuilt {
		fields := make([]string, 0, len(builder.fieldData))
		for k, v := range builder.fieldData {
			fields = append(fields, k+" = ?")
			builder.data = append(builder.data, v)
		}

		sql := fmt.Sprintf("update %s set %s", builder.table, strings.Join(fields, ", "))
		if builder.where != nil {
			where, whereData := builder.where.Build()
			sql = fmt.Sprintf("%s where %s", sql, where)
			builder.data = append(builder.data, whereData...)
		}
		builder.sql = sql
		builder.isBuilt = true
	}
	return builder.sql, builder.data
}
