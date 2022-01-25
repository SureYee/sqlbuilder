package sqlbuilder

func (builder *Builder) Limit(offset int, limit int) *Builder {
	builder.limit = limit
	builder.offset = offset
	return builder
}
