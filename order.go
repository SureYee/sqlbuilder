package sqlbuilder

func (builder *Builder) Order(column, sort string) *Builder {
	builder.order = append(builder.order, column+" "+sort)
	return builder
}
