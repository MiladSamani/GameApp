package main

import (
	"fmt"
	"gameAppProject/entity"
	"gameAppProject/repository/mysql"
)

func main() {

}

func testUserMysqlRepo() {
	mysqlRepo := mysql.New()
	createdUser, err := mysqlRepo.Register(entity.User{ID: 0, PhoneNumber: "093", Name: "Milad Samani"})
	if err != nil {
		fmt.Println("can't created register user", err)
	} else {
		fmt.Println("created user : ", createdUser)
	}
	isUnique, err := mysqlRepo.IsPhoneNumberUnique(createdUser.PhoneNumber + "23")
	if err != nil {
		fmt.Println("unique err", err)
	}

	fmt.Println("isUnique", isUnique)
}
