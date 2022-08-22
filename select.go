package sqlbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SelectBuilder 查询语句构建器
type SelectBuilder struct {
	where   WhereInterface
	having  WhereInterface
	isBuilt bool
	limit   int
	offset  int
	sql     string
	table   string
	join    []*Join
	fields  []string
	order   []string
	groupBy []string
	locker  Locker
	data    []interface{}
}

func (builder *SelectBuilder) Build() (string, []interface{}) {
	if !builder.isBuilt {
		sql := fmt.Sprintf("select %s from %s", strings.Join(builder.fields, ", "), builder.table)
		// 构建join
		if len(builder.join) > 0 {
			for _, j := range builder.join {
				join, joinData := j.Build()
				sql = sql + " " + join
				builder.data = append(builder.data, joinData...)
			}
		}

		// 构建where语句
		if builder.where != nil {
			where, whereData := builder.where.Build()
			if where != "" {
				sql = fmt.Sprintf("%s where %s", sql, where)
				builder.data = append(builder.data, whereData...)
			}
		}
		// 构建group by
		if len(builder.groupBy) > 0 {
			sql = fmt.Sprintf("%s group by %s", sql, strings.Join(builder.groupBy, ", "))
		}
		// 构建having
		if builder.having != nil {
			having, havingData := builder.having.Build()
			sql = fmt.Sprintf("%s having %s", sql, having)
			builder.data = append(builder.data, havingData...)
		}

		if len(builder.order) > 0 {
			sql = fmt.Sprintf("%s order by %s", sql, strings.Join(builder.order, ", "))
		}

		if builder.limit > 0 || builder.offset > 0 {
			sql = fmt.Sprintf("%s limit %d, %d", sql, builder.offset, builder.limit)
		}

		if builder.locker != nil {
			lockSql, _ := builder.locker.Build()
			sql = fmt.Sprintf("%s %s", sql, lockSql)
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

func (builder *SelectBuilder) getWhere() WhereInterface {
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

// Where some_column = value
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

func (builder *SelectBuilder) OrderBy(column, sort string) *SelectBuilder {
	sort = strings.ToLower(sort)
	if sort != "desc" && sort != "asc" {
		// 忽略错误的排序规则
		return builder
	}
	builder.order = append(builder.order, column+" "+sort)
	return builder
}

func (builder *SelectBuilder) GroupBy(column ...string) *SelectBuilder {
	builder.groupBy = append(builder.groupBy, column...)
	return builder
}

func (builder *SelectBuilder) Having(having BuilderFunc) *SelectBuilder {
	builder.having = &WhereBuilder{}
	builder.having.WhereFunc(having)
	return builder
}

func (builder *SelectBuilder) Limit(offset, limit int) *SelectBuilder {
	builder.limit = limit
	builder.offset = offset
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
		case reflect.Bool:
			return "1"
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

func (builder *SelectBuilder) LeftJoin(table string, on WhereInterface) *SelectBuilder {
	builder.join = append(builder.join, &Join{
		link:  "left",
		table: table,
		on:    on,
	})
	return builder
}

func (builder *SelectBuilder) RightJoin(table string, on WhereInterface) *SelectBuilder {
	builder.join = append(builder.join, &Join{
		link:  "right",
		table: table,
		on:    on,
	})
	return builder
}

func (builder *SelectBuilder) InnerJoin(table string, on WhereInterface) *SelectBuilder {
	builder.join = append(builder.join, &Join{
		link:  "inner",
		table: table,
		on:    on,
	})
	return builder
}

func (builder *SelectBuilder) LockForUpdate() *SelectBuilder {
	builder.locker = new(UpdateLocker)
	return builder
}

func (builder *SelectBuilder) LockShareMode() *SelectBuilder {
	builder.locker = new(ShareLocker)
	return builder
}
