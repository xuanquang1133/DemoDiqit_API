package main

import (
	"fmt"
	"demodiqit_api/config"
	"demodiqit_api/models"
)

func main() {
	cfg := config.LoadConfig()
	config.ConnectDB(cfg)

	var admin models.User
	if err := config.DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		fmt.Println("Admin not found")
		return
	}

	fmt.Println("Admin password hash in DB:", admin.Password)
}
