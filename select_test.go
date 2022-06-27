package sqlbuilder_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sureyee/sqlbuilder"
)

type user struct {
	status   int8
	username string
	mobile   string
	id       int
	age      int
}

func TestSelect(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "age", "status"}).
		AddRow(1, "zhangsan", 10, 1).
		AddRow(2, "lisi", 11, 1)

	query := "select * from users"
	mock.ExpectQuery(query).WillReturnRows(rows)

	buildQuery, builderData := sqlbuilder.Select("*").From("users").Build()
	row := db.QueryRow(buildQuery, builderData...)
	var data user
	if err := row.Scan(&data.id, &data.username, &data.age, &data.status); err != nil {
		t.Errorf("row.Scan error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWhere(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "age", "status"}).
		AddRow(1, "zhangsan", 10, 1).
		AddRow(2, "lisi", 11, 1)
	sql := "select * from users where id = ?"
	mock.ExpectQuery(sql).WillReturnRows(rows).WithArgs(2)
	builderSql, builderData := sqlbuilder.Select("*").From("users").Where("id", 2).Build()
	row := db.QueryRow(builderSql, builderData...)
	var data user
	if err := row.Scan(&data.id, &data.username, &data.age, &data.status); err != nil {
		t.Errorf("row.Scan error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMultiWhere(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "age", "status", "mobile"}).
		AddRow(1, "zhangsan", 10, 1, "13111112222").
		AddRow(2, "lisi", 11, 1, "13111111111")
	sql := "select * from users where username = ? and mobile = ?"
	mock.ExpectQuery(sql).WithArgs("zhangsan", "13111111111").WillReturnRows(rows)

	builderSql, builderData := sqlbuilder.Select("*").
		From("users").
		Where("username", "zhangsan").
		Where("mobile", "13111111111").
		Build()
	row := db.QueryRow(builderSql, builderData...)
	var data user
	if err := row.Scan(&data.id, &data.username, &data.age, &data.status, &data.mobile); err != nil {
		t.Errorf("row.Scan error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWhereIn(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "age", "status", "mobile"}).
		AddRow(1, "zhangsan", 10, 1, "13111112222").
		AddRow(2, "lisi", 11, 1, "13111111111")
	sql := "select * from users where id in (?, ?, ?, ?)"
	mock.ExpectQuery(sql).WithArgs(1, 2, 3, 4).WillReturnRows(rows)
	builderSql, builderData := sqlbuilder.Select("*").From("users").WhereIn("id", []int32{1, 2, 3, 4}).Build()
	row := db.QueryRow(builderSql, builderData...)
	var data user
	if err := row.Scan(&data.id, &data.username, &data.age, &data.status, &data.mobile); err != nil {
		t.Errorf("row.Scan error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWhereOperate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "age", "status", "mobile"}).
		AddRow(1, "zhangsan", 15, 1, "13111112222").
		AddRow(2, "lisi", 11, 1, "13111111111")
	sql := "select * from users where age < ?"
	mock.ExpectQuery(sql).WithArgs(12).WillReturnRows(rows)
	builderSql, builderData := sqlbuilder.Select("*").From("users").WhereOperate("age", "<", 12).Build()
	row := db.QueryRow(builderSql, builderData...)
	var data user
	if err := row.Scan(&data.id, &data.username, &data.age, &data.status, &data.mobile); err != nil {
		t.Errorf("row.Scan error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWhereLike(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "age", "status", "mobile"}).
		AddRow(1, "zhangsan", 15, 1, "13111112222").
		AddRow(2, "lisi", 11, 1, "13111111111")
	sql := "select * from users where username like ?"
	mock.ExpectQuery(sql).WithArgs("%zhang%").WillReturnRows(rows)

	builderSql, builderData := sqlbuilder.Select("*").From("users").WhereLike("username", "%zhang%").Build()
	row := db.QueryRow(builderSql, builderData...)
	var data user
	if err := row.Scan(&data.id, &data.username, &data.age, &data.status, &data.mobile); err != nil {
		t.Errorf("row.Scan error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
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

func TestGroupBy(t *testing.T) {
	sql := "select * from users group by age"
	builderSql := sqlbuilder.Select("*").From("users").GroupBy("age").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestHaving(t *testing.T) {
	sql := "select count(1) from users group by age having age > 10"
	builderSql := sqlbuilder.Select("count(1)").From("users").GroupBy("age").Having(func() sqlbuilder.Builder {
		return sqlbuilder.WhereOperate("age", ">", 10)
	}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestLimit(t *testing.T) {
	sql := "select * from users limit 0, 10"
	builderSql := sqlbuilder.Select("*").From("users").Limit(0, 10).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestOrderBy(t *testing.T) {
	sql := "select * from users order by create_time asc"
	builderSql := sqlbuilder.Select("*").From("users").OrderBy("create_time", "asc").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestOrderByDesc(t *testing.T) {
	sql := "select * from users order by create_time desc"
	builderSql := sqlbuilder.Select("*").From("users").OrderBy("create_time", "desc").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestErrorOrderBy(t *testing.T) {
	sql := "select * from users"
	builderSql := sqlbuilder.Select("*").From("users").OrderBy("create_time", "aaaa").String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestChildSelect(t *testing.T) {
	sql := "select * from users where id in (select user_id from books where is_publish = 1)"
	builderSql := sqlbuilder.Select("*").From("users").WhereIn("id", func() sqlbuilder.Builder {
		return sqlbuilder.Select("user_id").From("books").Where("is_publish", 1)
	}).String()

	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestWhereFunc(t *testing.T) {
	sql := "select * from users where gender = \"F\" and (age < 10 or age > 30)"
	builderSql := sqlbuilder.Select("*").From("users").Where("gender", "F").WhereFunc(func() sqlbuilder.Builder {
		return sqlbuilder.WhereOperate("age", "<", 10).OrWhereOperate("age", ">", 30)
	}).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestLeftJoin(t *testing.T) {
	sql := "select * from users left join books on books.user_id = users.id"
	builderSql := sqlbuilder.Select("*").From("users").LeftJoin(
		"books",
		sqlbuilder.WhereColumn("books.user_id", "users.id"),
	).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestRightJoin(t *testing.T) {
	sql := "select * from users right join books on books.user_id = users.id"
	builderSql := sqlbuilder.Select("*").From("users").RightJoin(
		"books",
		sqlbuilder.WhereColumn("books.user_id", "users.id"),
	).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestInnerJoin(t *testing.T) {
	sql := "select * from users inner join books on books.user_id = users.id"
	builderSql := sqlbuilder.Select("*").From("users").InnerJoin(
		"books",
		sqlbuilder.WhereColumn("books.user_id", "users.id"),
	).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}

func TestSelectBool(t *testing.T) {
	sql := "select * from users where status = 1"
	builderSql := sqlbuilder.Select("*").From("users").Where("status", true).String()
	if sql != builderSql {
		t.Errorf("expected:`%v`, got:`%v`", sql, builderSql)
	}
}
