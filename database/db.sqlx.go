package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var SqlxDB *sqlx.DB

func (d DB) InitSqlx() {
	db, err := sqlx.Open("mysql", "root:MlKX6KqF@tcp(127.0.0.1:3306)/cartoon")

	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	SqlxDB = db
	//defer db.Close()
}
