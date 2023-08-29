package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Auth struct {
		ExpiryDuration     int64
		ExpiryLongDuration int64
		SecretKey          string
	}
	RestServer struct {
		Address string
		Port    int
	}
	Database struct {
		Driver   string
		Host     string
		Port     int
		Username string
		Password string
		DBName   string
	}
}

var AppConfig Config

func (c *Config) GetDatabaseConnectionString() string {
	switch c.Database.Driver {
	case "mysql":
		// mysql://user:password@tcp(host:port)/dbname?query
		return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Database.Username,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.DBName,
		)
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?query&sslmode=disable",
			c.Database.Username,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.DBName,
		)
	case "sqlite3":
		return c.Database.DBName // SQLite3 might only need the path to the database file
	default:
		log.Fatalf("Unsupported database driver: %s", c.Database.Driver)
		return ""
	}
}

func InitConfig() {
	// Set the file name of the configurations file
	viper.SetConfigName("config.yaml")

	// Set the path to look for the configurations file
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")          // Optionally look for config in the working directory
	viper.AddConfigPath("config/")    // Path to look for the config file
	viper.AddConfigPath("/etc/wms/")  // Path to look for the config file
	viper.AddConfigPath("$HOME/.wms") // Call multiple times to add many search paths

	// Enable environment variable override
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Unmarshal the config into the AppConfig struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
