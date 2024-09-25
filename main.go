package main

import (
	"fmt"

	"apidemo/dbconfig"
)

func main() {
	err := dbconfig.DBConnect()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	println("Hello, World!")
}
