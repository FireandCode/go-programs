package main

import (
	"fmt"
	"sync"
)

type DatabaseConnection struct{}

func (db *DatabaseConnection) Query(sql string) {
	fmt.Println("Executing query:", sql)
}

var (
	DB   *DatabaseConnection
	once sync.Once
)

func Connection() *DatabaseConnection {
	once.Do(func() {
		fmt.Println("Creating a new DB connection")
		DB = &DatabaseConnection{}
	})
	return DB
}

func main() {
	db1 := Connection()
	db1.Query("SELECT * FROM users")

	db2 := Connection()
	db2.Query("SELECT * FROM products")

	fmt.Println("Same instance?", db1 == db2) // true
}
