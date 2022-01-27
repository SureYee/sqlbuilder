package sqlbuilder

// Builder Builder接口
// 在每个实现builder接口的结构体中，执行build方法，将会返回构建的sql和data
type Builder interface {
	Build() (string, []interface{})
}

type BuilderFunc func() Builder

type Column string
