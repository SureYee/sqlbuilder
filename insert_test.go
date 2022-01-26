package sqlbuilder_test

import (
	"testing"

	"github.com/sureyee/sqlbuilder"
)

func TestInsert(t *testing.T) {
	sql := "insert into users values (\"zhangsan\", 10)"
	builderSql := sqlbuilder.Insert("users").Values("zhangsan", 10).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestInsertFields(t *testing.T) {
	sql1 := "insert into users (username, age) values (\"zhangsan\", 10)"
	sql2 := "insert into users (age, username) values (10, \"zhangsan\")"
	builderSql := sqlbuilder.Insert("users").Fields("username", "age").Values("zhangsan", 10).String()
	if sql1 == builderSql {
		return
	}

	if sql2 == builderSql {
		return
	}

	t.Errorf("expected:`%v` or `%v`, got:`%v`", sql1, sql2, builderSql)
}

func TestInsertMap(t *testing.T) {
	sql := "insert into users (username, age) values (\"zhangsan\", 10)"
	builderSql := sqlbuilder.Insert("users").Map(map[string]interface{}{
		"username": "zhangsan",
		"age":      10,
	}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}
