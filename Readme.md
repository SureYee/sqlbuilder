# sqlbuilder

## 安装

```shell
go get -u github.com/sureyee/sqlbuilder@v0.1.0-alpha
```

## 使用

### 查询

1. 条件查询

```go
sql, data := sqlbuilder.Select("*").From("users").Where("id", 10).Build()
// select * from users where id = ?
// [10]

sql, data := sqlbuilder.Select("*").From("users").WhereOperate("age", "<", 10).OrWhereOperate("age", ">", 50).Build()
```



