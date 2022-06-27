# sqlbuilder

## 安装

```shell
go get -u github.com/sureyee/sqlbuilder@v0.1.0-alpha
```

## 使用

### 查询

1. 条件查询

```go
// 普通查询
sql, data := sqlbuilder.Select("*").From("users").Where("id", 10).Build()
// select * from users where id = ?
// [10]

// or语句查询
sql, data := sqlbuilder.Select("*").From("users").WhereOperate("age", "<", 10).OrWhereOperate("age", ">", 50).Build()
// select * from users where age < ? or age > ?  
// [10 50]

// ()查询, 使用WhereFunc闭包进行查询回将查询语句用()包裹
sql, data := sqlbuilder.Select("*").From("users").Where("gender", "F").WhereFunc(func() sqlbuilder.Builder {
		return sqlbuilder.WhereOperate("age", "<", 10).OrWhereOperate("age", ">", 30)
	}).Build()
//select * from users where gender = ? and (age < ? or age > ?) 
//[F 10 30]

// in 子查询语句
sql, data := sqlbuilder.Select("*").From("users").WhereIn("id", func() sqlbuilder.Builder {
		return sqlbuilder.Select("user_id").From("books").Where("is_publish", 1)
	}).Build()


```

> 更多用法查看测试文件



