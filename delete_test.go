package sqlbuilder_test

import (
	"testing"

	"sureyee.com/sqlbuilder"
)

func TestDelete(t *testing.T) {
	sql := "delete from users"
	builderSql := sqlbuilder.Delete("users").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestDeleteWhere(t *testing.T) {
	sql := "delete from users where id = 10"
	builderSql := sqlbuilder.Delete("users").Where("id", 10).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestDeleteWhereIn(t *testing.T) {
	sql := "delete from users where id in (1, 2, 3)"
	builderSql := sqlbuilder.Delete("users").WhereIn("id", []int{1, 2, 3}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}
