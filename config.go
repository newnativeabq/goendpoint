package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configurations Exported
type Configurations struct {
	Server       ServerConfigurations
	Database     DatabaseConfigurations
	EXAMPLE_PATH string
	EXAMPLE_VAR  string
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port    int
	Address string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBName     string
	DBUser     string
	DBPassword string
}

// BuildConfigurations generates configurations structure from .yml file
func BuildConfigurations(filename string, filetype string) Configurations {
	// Set the file name of the configuration file
	viper.SetConfigName(filename)

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType(filetype)
	var configuration Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return configuration
}