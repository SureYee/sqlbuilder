package sqlbuilder

import (
	"fmt"
	"strings"
)

func (builder *Builder) Group(fileds ...string) *Builder {
	builder.groupBy = "group by " + strings.Join(fileds, ",")
	return builder
}

func (builder *Builder) Having(column, operate string, value interface{}) *Builder {
	return builder.where(column, operate, value)
}

func (builder *Builder) buildHaving() (string, []interface{}) {
	var having string
	data := make([]interface{}, 0)
	for _, v := range builder.having {
		switch v.value.(type) {
		case *Builder:
			b := v.value.(*Builder)
			if v.expression != "" {
				//TODO where id = (sql)
			} else {
				w, d := b.buildWhere()
				data = append(data, d...)
				if having == "" {
					having += fmt.Sprintf("(%s)", w)
				} else {
					having += fmt.Sprintf(" and (%s)", w)
				}
			}
		case string, int:
			if having == "" {
				having += v.expression
			} else {
				having += " and " + v.expression
			}
			data = append(data, v.value)
		case []interface{}:
			if having == "" {
				having += v.expression
			} else {
				having += " and " + v.expression
			}
			value := v.value.([]interface{})
			data = append(data, value...)
		}
	}
	return having, data
}
