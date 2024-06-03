/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"courier_service/src/cmd"
	"fmt"

	"courier_service/config"
)

func main() {
	appConfig := config.NewConfig()
	err := appConfig.LoadConfig("config/app_config.json")
	if err != nil {
		fmt.Errorf("Error while loading configs:%s\n", err)
		return
	}
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
