package sqlbuilder

import "fmt"

func (builder *Builder) LeftJoin(table, on string) *Builder {
	builder.join = append(builder.join, fmt.Sprintf("left join %s on %s", table, on))
	return builder
}

func (builder *Builder) RightJoin(table, on string) *Builder {
	builder.join = append(builder.join, fmt.Sprintf("right join %s on %s", table, on))
	return builder
}

func (builder *Builder) InnerJoin(table, on string) *Builder {
	builder.join = append(builder.join, fmt.Sprintf("inner join %s on %s", table, on))
	return builder
}
