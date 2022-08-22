package sqlbuilder

// Builder Builder接口
// 在每个实现builder接口的结构体中，执行build方法，将会返回构建的sql和data
type Builder interface {
	Build() (string, []interface{})
}

type BuilderFunc func() Builder

type Column string

type RawExpr struct {
	expr string
	data []interface{}
}

func Raw(expr string, data ...interface{}) *RawExpr {
	return &RawExpr{
		expr: expr,
		data: data,
	}
}

func (r *RawExpr) Build() (string, []interface{}) {
	return r.expr, r.data
}
