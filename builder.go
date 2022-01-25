package sqlbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

type Builder struct {
	isBuilt   bool
	statement string
	table     string
	fields    string
	sql       string
	groupBy   string
	limit     int
	offset    int
	wh        []*where
	orWh      []*where
	join      []string
	order     []string
	having    []*where
	data      []interface{}
}

func (builder *Builder) Select(fields ...string) *Builder {
	builder.statement = "select"
	builder.fields = strings.Join(fields, ", ")
	return builder
}

func (builder *Builder) Update(table string, data map[string]interface{}) *Builder {
	builder.statement = "update"
	builder.table = table
	fields := make([]string, 0, len(data))
	for k, v := range data {
		fields = append(fields, k+"=?")
		builder.data = append(builder.data, v)
	}
	builder.fields = strings.Join(fields, ",")

	return builder
}

func (builder *Builder) Delete(table string) *Builder {
	builder.statement = "delete"
	builder.table = table
	return builder
}

func (builder *Builder) From(table string) *Builder {
	builder.table = table
	return builder
}

func (builder *Builder) Perpare() (string, []interface{}) {
	if !builder.isBuilt {
		switch builder.statement {
		case "select":
			builder.sql, builder.data = builder.buildSelect()
		case "update":
			builder.sql, builder.data = builder.buildUpdate()
		case "delete":
			builder.sql, builder.data = builder.buildDelete()
		default:
			builder.sql, builder.data = builder.buildSelect()
		}
	}
	return builder.sql, builder.data
}

func (builder *Builder) String() string {
	sql, data := builder.Perpare()
	index := 0
	newSql := make([]rune, 0, len(data))
	getData := func(data []interface{}, index int) string {
		if index > (len(data) - 1) {
			return ""
		}

		datum := data[index]
		switch datum := datum.(type) {
		case int:
			return strconv.Itoa(datum)
		case string:
			return "\"" + strings.ReplaceAll(strings.ReplaceAll(datum, "\\", "\\\\"), "\"", "\\\"") + "\""
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

func (builder *Builder) buildSelect() (string, []interface{}) {
	sql := fmt.Sprintf("select %s from %s", builder.fields, builder.table)

	if len(builder.join) > 0 {
		for _, s := range builder.join {
			sql += " " + s
		}
	}

	data := make([]interface{}, 0)
	if len(builder.wh) > 0 {
		where, whereData := builder.buildWhere()
		data = append(data, whereData...)
		sql = fmt.Sprintf("%s where %s", sql, where)
		if len(builder.orWh) > 0 {
			where, whereData := builder.buildOrWhere()
			data = append(data, whereData...)
			sql = fmt.Sprintf("%s or %s", sql, where)
		}
	} else {
		if len(builder.orWh) > 0 {
			where, whereData := builder.buildOrWhere()
			data = append(data, whereData...)
			sql = fmt.Sprintf("%s where %s", sql, where)
		}
	}

	if len(builder.having) > 0 {
		having, havingData := builder.buildHaving()
		data = append(data, havingData...)
		sql = fmt.Sprintf("%s having %s", sql, having)
	}

	if len(builder.order) > 0 {
		sql = fmt.Sprintf("%s order by %s", sql, strings.Join(builder.order, ","))
	}

	if builder.limit != 0 {
		sql = fmt.Sprintf("%s limit %d, %d", sql, builder.offset, builder.limit)
	}

	return sql, data
}

func (builder *Builder) buildUpdate() (string, []interface{}) {
	sql := fmt.Sprintf("update %s set %s", builder.table, builder.fields)

	if len(builder.wh) > 0 {
		where, whereData := builder.buildWhere()
		builder.data = append(builder.data, whereData...)
		sql = fmt.Sprintf("%s where %s", sql, where)
	}

	if len(builder.orWh) > 0 {
		where, whereData := builder.buildOrWhere()
		builder.data = append(builder.data, whereData...)
		sql = fmt.Sprintf("%s or (%s)", sql, where)
	}

	return sql, builder.data
}

func (builder *Builder) buildDelete() (string, []interface{}) {
	sql := fmt.Sprintf("delete from %s", builder.table)

	if len(builder.wh) > 0 {
		where, whereData := builder.buildWhere()
		builder.data = append(builder.data, whereData...)
		sql = fmt.Sprintf("%s where %s", sql, where)
	}

	if len(builder.orWh) > 0 {
		where, whereData := builder.buildOrWhere()
		builder.data = append(builder.data, whereData...)
		sql = fmt.Sprintf("%s or (%s)", sql, where)
	}
	return sql, builder.data
}

func NewBuidler() *Builder {
	return &Builder{}
}
