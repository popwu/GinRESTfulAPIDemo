package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// 从参数获取原始密码，加密后输出
func main() {
	password := os.Args[1]
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating password:", err)
		return
	}
	fmt.Println(string(hashedPassword))
}
