package sqlbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type DeleteBuilder struct {
	isBuilt bool
	sql     string
	table   string
	where   *WhereBuilder
	data    []interface{}
}

func Delete(table string) *DeleteBuilder {
	return &DeleteBuilder{
		table: table,
	}
}

func (builder *DeleteBuilder) getWhere() where {
	if builder.where == nil {
		builder.where = &WhereBuilder{}
	}
	return builder.where
}

// WhereOperate where column > value
// 可以指定操作符的where语句
func (builder *DeleteBuilder) WhereOperate(column, operate string, value interface{}) *DeleteBuilder {
	builder.getWhere().WhereOperate(column, operate, value)
	return builder
}

// Where where some_column = value
// 普通的where语句
func (builder *DeleteBuilder) Where(column string, value interface{}) *DeleteBuilder {
	builder.getWhere().Where(column, value)
	return builder
}

func (builder *DeleteBuilder) WhereIn(column string, value interface{}) *DeleteBuilder {
	builder.getWhere().WhereIn(column, value)
	return builder
}

func (builder *DeleteBuilder) WhereBetween(column string, min, max interface{}) *DeleteBuilder {
	builder.getWhere().WhereBetween(column, min, max)
	return builder
}

func (builder *DeleteBuilder) WhereLike(column string, value interface{}) *DeleteBuilder {
	builder.getWhere().WhereLike(column, value)
	return builder
}

func (builder *DeleteBuilder) WhereNull(column string) *DeleteBuilder {
	builder.getWhere().WhereNull(column)
	return builder
}

func (builder *DeleteBuilder) WhereNotNull(column string) *DeleteBuilder {
	builder.getWhere().WhereNotNull(column)
	return builder
}

func (builder *DeleteBuilder) WhereFunc(f BuilderFunc) *DeleteBuilder {
	builder.getWhere().WhereFunc(f)
	return builder
}

func (builder *DeleteBuilder) OrWhereOperate(column, operate string, value interface{}) *DeleteBuilder {
	builder.getWhere().OrWhereOperate(column, operate, value)
	return builder
}

func (builder *DeleteBuilder) OrWhere(column string, value interface{}) *DeleteBuilder {
	builder.getWhere().OrWhere(column, value)
	return builder
}

func (builder *DeleteBuilder) OrWhereFunc(f BuilderFunc) *DeleteBuilder {
	builder.getWhere().OrWhereFunc(f)
	return builder
}

func (builder *DeleteBuilder) String() string {
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

func (builder *DeleteBuilder) Build() (string, []interface{}) {
	if !builder.isBuilt {
		sql := fmt.Sprintf("delete from %s", builder.table)
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
