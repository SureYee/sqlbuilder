package sqlbuilder_test

import (
	"testing"

	"sureyee.com/sqlbuilder"
)

func TestWhere(t *testing.T) {
	sql := "select * from users where id = 10"
	builderSql := sqlbuilder.Select("*").From("users").Where("id", 10).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestMultiWhere(t *testing.T) {
	sql := "select * from users where username = \"zhangsan\" and mobile = \"13111111111\""
	builderSql := sqlbuilder.Select("*").From("users").Where("username", "zhangsan").Where("mobile", "13111111111").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestWhereIn(t *testing.T) {
	sql := "select * from users where id in (1, 2, 3, 4)"
	builderSql := sqlbuilder.Select("*").From("users").WhereIn("id", []int32{1, 2, 3, 4}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestWhereOperate(t *testing.T) {
	sql := "select * from users where age < 10"
	builderSql := sqlbuilder.Select("*").From("users").WhereOperate("age", "<", 10).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestWhereLike(t *testing.T) {
	sql := "select * from users where username like \"%zhang%\""
	builderSql := sqlbuilder.Select("*").From("users").WhereLike("username", "%zhang%").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestWhereNull(t *testing.T) {
	sql := "select * from users where is_delete is null"
	builderSql := sqlbuilder.Select("*").From("users").WhereNull("is_delete").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestWhereNotNull(t *testing.T) {
	sql := "select * from users where is_delete is not null"
	builderSql := sqlbuilder.Select("*").From("users").WhereNotNull("is_delete").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestOrWhere(t *testing.T) {
	sql := "select * from users where username = \"zhangsan\" or username = \"lisi\""
	builderSql := sqlbuilder.Select("*").From("users").Where("username", "zhangsan").OrWhere("username", "lisi").String()

	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestOrWhereFunc(t *testing.T) {
	sql := "select * from users where username = \"zhangsan\" or (mobile = \"13111111111\" and name = \"张三\")"
	builderSql := sqlbuilder.Select("*").From("users").Where("username", "zhangsan").OrWhereFunc(func() sqlbuilder.Builder {
		return sqlbuilder.Where("mobile", "13111111111").Where("name", "张三")
	}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestAllWhere(t *testing.T) {
	sql := "select * from users where age < 10 and id in (1, 2, 3, 4) and create_time between \"2020-01\" and \"2020-02\" and is_delete is not null or (username like \"zhangsan\" and (age < 20 or age > 30))"
	builderSql := sqlbuilder.Select("*").From("users").WhereOperate("age", "<", 10).WhereIn("id", []int{1, 2, 3, 4}).WhereBetween("create_time", "2020-01", "2020-02").WhereNotNull("is_delete").OrWhereFunc(func() sqlbuilder.Builder {
		return sqlbuilder.WhereLike("username", "zhangsan").WhereFunc(func() sqlbuilder.Builder {
			return sqlbuilder.WhereOperate("age", "<", 20).OrWhereOperate("age", ">", 30)
		})
	}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}
