package sqlbuilder_test

import (
	"fmt"
	"testing"

	"sureyee.com/sqlbuilder"
)

func TestSelect(t *testing.T) {
	sql := "select * from users"
	builder := sqlbuilder.NewBuidler().Select("*").From("users")

	if sql != builder.String() {
		t.Errorf("expected:%v, got:%v", sql, builder.String())
	}

}

func TestUpdate(t *testing.T) {
	builder := sqlbuilder.NewBuidler().Update("users", map[string]interface{}{
		"username": "zhangsan",
		"password": "lsdajflksjf",
	}).Where("id", 10).OrWhere("username", "lisi")
	prepare, data := builder.Perpare()
	fmt.Printf("%s %+v\n", prepare, data)
	fmt.Println(builder.String())
}
