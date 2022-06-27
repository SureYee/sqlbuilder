package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"
)

type WhereInterface interface {
	Where(string, interface{}) WhereInterface
	WhereColumn(string, string) WhereInterface
	WhereColumnOperate(string, string, string) WhereInterface
	WhereBetween(string, interface{}, interface{}) WhereInterface
	WhereLike(string, interface{}) WhereInterface
	WhereIn(string, interface{}) WhereInterface
	WhereNull(string) WhereInterface
	WhereNotNull(string) WhereInterface
	WhereOperate(string, string, interface{}) WhereInterface
	WhereFunc(BuilderFunc) WhereInterface
	OrWhere(string, interface{}) WhereInterface
	OrWhereOperate(string, string, interface{}) WhereInterface
	OrWhereFunc(BuilderFunc) WhereInterface
	Build() (string, []interface{})
}

type WhereBuilder struct {
	wh   []*whereStat
	orWh []*whereStat
}

type whereStat struct {
	column  Column
	operate string
	value   interface{}
}

func Where(column string, value interface{}) *WhereBuilder {
	builder := &WhereBuilder{}
	builder.Where(column, value)
	return builder
}

func WhereLike(column string, value interface{}) *WhereBuilder {
	builder := &WhereBuilder{}
	builder.WhereLike(column, value)
	return builder
}

func WhereNull(column string) *WhereBuilder {
	builder := &WhereBuilder{}
	builder.WhereNull(column)
	return builder
}

func WhereOperate(column, operate string, value interface{}) *WhereBuilder {
	builder := &WhereBuilder{}
	builder.WhereOperate(column, operate, value)
	return builder
}

func WhereColumn(column1, column2 string) *WhereBuilder {
	builder := &WhereBuilder{}
	builder.WhereColumn(column1, column2)
	return builder
}

func WhereColumnOperate(column1, operate, column2 string) *WhereBuilder {
	builder := &WhereBuilder{}
	builder.WhereColumnOperate(column1, operate, column2)
	return builder
}

func (stat *whereStat) Build() (string, []interface{}) {
	switch stat.operate {
	case "in":
		return stat.buildIn()
	case "between":
		return stat.buildBetween()
	case "build":
		return stat.buildSql()
	case "is":
		return stat.buildIs()
	case "not":
		return stat.buildNot()
	default:
		switch f := stat.value.(type) {
		case func() Builder:
			sql, data := f().Build()
			return fmt.Sprintf("%s %s (%s)", stat.column, stat.operate, sql), data
		case BuilderFunc:
			sql, data := f().Build()
			return fmt.Sprintf("%s %s (%s)", stat.column, stat.operate, sql), data
		case Column:
			return fmt.Sprintf("%s %s %s", stat.column, stat.operate, stat.value), nil
		}
		return fmt.Sprintf("%s %s ?", stat.column, stat.operate), []interface{}{stat.value}
	}
}

func (stat *whereStat) buildIs() (string, []interface{}) {
	sql := fmt.Sprintf("%s is null", stat.column)
	return sql, nil
}

func (stat *whereStat) buildNot() (string, []interface{}) {
	sql := fmt.Sprintf("%s is not null", stat.column)
	return sql, nil
}

func (stat *whereStat) buildSql() (string, []interface{}) {
	if v, ok := stat.value.(Builder); ok {
		if w, ok := v.(*WhereBuilder); ok {
			sql, data := w.Build()
			if len(w.wh)+len(w.orWh) > 1 {
				return "(" + sql + ")", data
			}
			return sql, data
		}
		return v.Build()
	}
	panic("where func value must Builder")
}

func (stat *whereStat) buildBetween() (string, []interface{}) {
	sql := fmt.Sprintf("%s between ? and ?", stat.column)

	if v, ok := stat.value.([]interface{}); ok {
		return sql, v
	}
	panic("between value must []interface{}")
}

func (stat *whereStat) buildIn() (string, []interface{}) {
	v := reflect.ValueOf(stat.value)
	switch v.Kind() {
	case reflect.Slice:
		data := make([]interface{}, v.Len())
		replace := make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			replace[i] = "?"
			switch v.Index(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				data[i] = int(v.Index(i).Int())
			case reflect.String:
				data[i] = v.Index(i).String()
			case reflect.Interface:
				switch inter := v.Index(i).Interface().(type) {
				case string, int, int8, int16, int32:
					data[i] = inter
				default:
					panic("build in value is invalid")
				}
			default:
				data[i] = v.Index(i).String()
			}
		}
		sql := fmt.Sprintf("%s in (%s)", stat.column, strings.Join(replace, ", "))
		return sql, data
	case reflect.Func:
		switch f := stat.value.(type) {
		case func() Builder:
			sql, data := f().Build()
			return fmt.Sprintf("%s in (%s)", stat.column, sql), data
		case BuilderFunc:
			sql, data := f().Build()
			return fmt.Sprintf("%s in (%s)", stat.column, sql), data
		}
		panic("where func must BuilderFunc")
	default:
		panic("where in value must slice")
	}

}

func (builder *WhereBuilder) Where(column string, value interface{}) WhereInterface {
	return builder.WhereOperate(column, "=", value)
}

func (builder *WhereBuilder) WhereColumn(column1, column2 string) WhereInterface {
	return builder.WhereOperate(column1, "=", Column(column2))
}

func (builder *WhereBuilder) WhereColumnOperate(column1, operate, column2 string) WhereInterface {
	return builder.WhereOperate(column1, operate, Column(column2))
}

func (builder *WhereBuilder) WhereBetween(column string, min, max interface{}) WhereInterface {
	return builder.WhereOperate(column, "between", []interface{}{min, max})
}

func (builder *WhereBuilder) WhereIn(column string, value interface{}) WhereInterface {
	return builder.WhereOperate(column, "in", value)
}

func (builder *WhereBuilder) WhereLike(column string, value interface{}) WhereInterface {
	return builder.WhereOperate(column, "like", value)
}

func (builder *WhereBuilder) WhereNull(column string) WhereInterface {
	return builder.WhereOperate(column, "is", nil)
}

func (builder *WhereBuilder) WhereNotNull(column string) WhereInterface {
	return builder.WhereOperate(column, "not", nil)
}

// WhereOperate where column > value
// 可以指定操作符的where语句
func (builder *WhereBuilder) WhereOperate(column, operate string, value interface{}) WhereInterface {
	builder.wh = append(builder.wh, &whereStat{
		column:  Column(column),
		operate: operate,
		value:   value,
	})
	return builder
}

func (builder *WhereBuilder) WhereFunc(f BuilderFunc) WhereInterface {
	builder.wh = append(builder.wh, &whereStat{
		operate: "build",
		value:   f(),
	})
	return builder
}

func (builder *WhereBuilder) OrWhere(column string, value interface{}) WhereInterface {
	return builder.OrWhereOperate(column, "=", value)
}

func (builder *WhereBuilder) OrWhereOperate(column, operate string, value interface{}) WhereInterface {
	builder.orWh = append(builder.orWh, &whereStat{
		column:  Column(column),
		operate: operate,
		value:   value,
	})
	return builder
}

func (builder *WhereBuilder) OrWhereFunc(f BuilderFunc) WhereInterface {
	builder.orWh = append(builder.orWh, &whereStat{
		operate: "build",
		value:   f(),
	})
	return builder
}

func (builder *WhereBuilder) Build() (string, []interface{}) {
	sql := ""
	data := make([]interface{}, 0)
	if len(builder.wh) > 0 {
		// builder where
		for _, v := range builder.wh {
			w, d := v.Build()
			if sql == "" {
				sql = w
				data = d
			} else {
				sql = fmt.Sprintf("%s and %s", sql, w)
				data = append(data, d...)
			}
		}

		if len(builder.orWh) > 0 {
			// build or where
			for _, v := range builder.orWh {
				w, d := v.Build()
				if sql == "" {
					sql = w
					data = d
				} else {
					sql = fmt.Sprintf("%s or %s", sql, w)
					data = append(data, d...)
				}
			}
		}
	}

	return sql, data
}
