package sqlbuilder_test

import (
	"testing"

	"sureyee.com/sqlbuilder"
)

func TestUpdate(t *testing.T) {
	sql := "update users set username = \"zhangsan\""
	builderSql := sqlbuilder.Update("users").Set("username", "zhangsan").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestUpdateColumns(t *testing.T) {
	sql1 := "update users set username = \"zhangsan\", age = 10"
	sql2 := "update users set age = 10, username = \"zhangsan\""
	builderSql := sqlbuilder.Update("users").Set("username", "zhangsan").Set("age", 10).String()

	if sql1 == builderSql {
		return
	}

	if sql2 == builderSql {
		return
	}

	t.Errorf("expected:`%v` or `%v`, got:`%v`", sql1, sql2, builderSql)
}

func TestUpateMap(t *testing.T) {
	sql1 := "update users set username = \"zhangsan\", age = 10"
	sql2 := "update users set age = 10, username = \"zhangsan\""
	builderSql := sqlbuilder.Update("users").Map(map[string]interface{}{
		"username": "zhangsan",
		"age":      10,
	}).String()

	if sql1 == builderSql {
		return
	}

	if sql2 == builderSql {
		return
	}

	t.Errorf("expected:`%v` or `%v`, got:`%v`", sql1, sql2, builderSql)
}

func TestUpdateWhere(t *testing.T) {
	sql := "update users set username = \"zhangsan\" where id = 1"
	builderSql := sqlbuilder.Update("users").Set("username", "zhangsan").Where("id", 1).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestUpdateWhereOperate(t *testing.T) {
	sql := "update users set age = 99 where age > 99"
	builderSql := sqlbuilder.Update("users").Set("age", 99).WhereOperate("age", ">", 99).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestUpdateWhereIn(t *testing.T) {
	sql := "update users set age = 99 where id in (1, 2, 3)"
	builderSql := sqlbuilder.Update("users").Set("age", 99).WhereIn("id", []int{1, 2, 3}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestUpdateWhereBetween(t *testing.T) {
	sql := "update users set age = 10 where id between 1 and 10"
	builderSql := sqlbuilder.Update("users").Set("age", 10).WhereBetween("id", 1, 10).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}
