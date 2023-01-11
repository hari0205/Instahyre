package main

import (
	ini "example.com/Instahyre/teleapi/init"
	model "example.com/Instahyre/teleapi/models"
)

func init() {
	ini.ConnectToDB()
}

func main() {
	ini.DB.AutoMigrate(&model.UserData{})
}
