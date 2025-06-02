package main

import (
	"fmt"
	"idm/inner/database"
	"idm/inner/role"
)

func main() {
	fmt.Println("Hello Go")
	db, _ := database.ConnectDb()
	service := role.NewService(role.NewRoleRepository(db))
	save, _ := service.FindAll()
	fmt.Println(save)
}
