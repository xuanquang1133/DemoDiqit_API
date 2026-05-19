package main

import (
	"fmt"
	"demodiqit_api/config"
	"demodiqit_api/models"
)

func main() {
	cfg := config.LoadConfig()
	config.ConnectDB(cfg)

	// Update existing users to have Status "Active" if empty
	if err := config.DB.Model(&models.User{}).Where("status = '' OR status IS NULL").Update("status", "Active").Error; err != nil {
		fmt.Println("Error updating user statuses:", err)
		return
	}

	fmt.Println("Successfully updated old user statuses to Active!")
}
