package sqlbuilder

import "fmt"

type Join struct {
	link  string
	table string
	on    where
}

func (builder *Join) Build() (string, []interface{}) {
	sql, data := builder.on.Build()
	return fmt.Sprintf("%s join %s on %s", builder.link, builder.table, sql), data
}
