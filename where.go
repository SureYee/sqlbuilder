package sqlbuilder

import (
	"fmt"
	"strings"
)

type where struct {
	expression string
	value      interface{}
}

type whereFunc func() *Builder

// Where where some_column = value
// 普通的where语句
func (builder *Builder) Where(column string, value interface{}) *Builder {
	return builder.where(column, "=", value)
}

// WhereOperate where column > value
// 可以指定操作符的where语句
func (builder *Builder) WhereOperate(column, operate string, value interface{}) *Builder {
	return builder.where(column, operate, value)
}

// WhereLike where like value
func (builder *Builder) WhereLike(column string, value interface{}) *Builder {
	return builder.where(column, "like", value)
}

// WhereFunc where (somewhere)
// 当需要使用()时，使用此方法
func (builder *Builder) WhereFunc(f whereFunc) *Builder {
	builder.wh = append(builder.wh, &where{
		value: f(),
	})

	return builder
}

// WhereBetween between min and max
// between 语句
func (builder *Builder) WhereBetween(column string, min interface{}, max interface{}) *Builder {
	builder.wh = append(builder.wh, &where{
		expression: fmt.Sprintf("%s between ? and ?", column),
		value:      []interface{}{min, max},
	})
	return builder
}

// WhereIn in (,,,)
// in () 语句
func (builder *Builder) WhereIn(column string, value []interface{}) *Builder {
	replace := make([]string, len(value))
	for i := 0; i < len(value); i++ {
		replace[i] = "?"
	}
	builder.wh = append(builder.wh, &where{
		expression: fmt.Sprintf("%s in (%s)", column, strings.Join(replace, ", ")),
		value:      value,
	})
	return builder
}

// where push item to builder.wh
func (builder *Builder) where(column, operate string, value interface{}) *Builder {
	builder.wh = append(builder.wh, &where{
		expression: fmt.Sprintf("%s %s ?", column, operate),
		value:      value,
	})
	return builder
}

// OrWhere or column = vlaue
// or 语句
func (builder *Builder) OrWhere(column string, value interface{}) *Builder {
	builder.orWh = append(builder.orWh, &where{
		expression: fmt.Sprintf("%s %s ?", column, "="),
		value:      value,
	})
	return builder
}

func (builder *Builder) OrWhereOperate(column, operate string, value interface{}) *Builder {
	builder.orWh = append(builder.orWh, &where{
		expression: fmt.Sprintf("%s %s ?", column, operate),
		value:      value,
	})
	fmt.Printf("%#v", builder)
	return builder
}

// OrWhereFunc or (some where sql)
// or () 语句带括号
func (builder *Builder) OrWhereFunc(f whereFunc) *Builder {
	builder.orWh = append(builder.orWh, &where{
		value: f(),
	})
	return builder
}

// buildWhere
// 构建where语句
func (builder *Builder) buildWhere() (string, []interface{}) {
	var where string
	data := make([]interface{}, 0)
	for _, v := range builder.wh {
		fmt.Println(v)
		switch v.value.(type) {
		case *Builder:
			b := v.value.(*Builder)
			if v.expression != "" {
				//TODO where id = (sql)
			} else {
				if len(b.wh) > 0 {
					w, d := b.buildWhere()
					data = append(data, d...)
					if where == "" {
						where += fmt.Sprintf("(%s)", w)
					} else {
						where += fmt.Sprintf(" and (%s)", w)
					}
					if len(b.orWh) > 0 {
						w, d := b.buildOrWhere()
						data = append(data, d...)
						where = fmt.Sprintf("%s or %s", where, w)
					}
				} else {
					if len(b.orWh) > 0 {
						w, d := builder.buildOrWhere()
						data = append(data, d...)
						where = fmt.Sprintf("%s where %s", where, w)
					}
				}

			}
		case string, int:
			if where == "" {
				where += v.expression
			} else {
				where += " and " + v.expression
			}
			data = append(data, v.value)
		case []interface{}:
			if where == "" {
				where += v.expression
			} else {
				where += " and " + v.expression
			}
			value := v.value.([]interface{})
			data = append(data, value...)
		}
	}
	return where, data
}

// buildOrWhere
// 构建or语句
func (builder *Builder) buildOrWhere() (string, []interface{}) {
	var where string
	data := make([]interface{}, 0)
	for _, v := range builder.orWh {
		switch v.value.(type) {
		case *Builder:
			b := v.value.(*Builder)
			if v.expression != "" {
				//TODO where id = (sql)
			} else {
				// orWhereFunc 需要加上()
				w, d := b.buildWhere()
				data = append(data, d...)
				if where == "" {
					where += fmt.Sprintf("(%s)", w)
				} else {
					where += fmt.Sprintf(" or (%s)", w)
				}
			}
		case string, int:
			if where == "" {
				where += v.expression
			} else {
				where += " or " + v.expression
			}
			data = append(data, v.value)
		case []interface{}:
			if where == "" {
				where += v.expression
			} else {
				where += " or " + v.expression
			}
			value := v.value.([]interface{})
			data = append(data, value...)
		}
	}
	return where, data
}
