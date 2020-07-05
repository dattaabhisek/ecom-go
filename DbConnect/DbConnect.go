package DbConnect

import (

	//"fmt"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
)

func DbConnect() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "bankura-2020"
    dbHost := "localhost"
    dbPort := "3306"
    dbName := "goecom"
    db, err := sql.Open(dbDriver, dbUser +":"+ dbPass +"@tcp("+ dbHost +":"+ dbPort +")/"+ dbName +"?charset=utf8")
    checkErr(err)
    return db
}
func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}