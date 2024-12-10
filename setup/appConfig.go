package setup

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type AppConfig struct {
	Host         string      `json:"host"`
	DatabaseFile string      `json:"databaseFile"`
	DefaultUser  DefaultUser `json:"defaultUser"`
	JWTSecret    string      `json:"jwtSecret"`
}

type DefaultUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ReadConfig() AppConfig {
	configPath := readConfigPath()

	configJson, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file in %s", configPath)
	}

	var config AppConfig
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		log.Fatalf("Failed to read config file in %s", configPath)
	}

	return config
}

func readConfigPath() string {
	configPathVar := flag.String("config", "config.json", "Path to json config for the server")
	flag.Parse()
	log.Printf("Loading config from \"%s\"", *configPathVar)
	return *configPathVar
}
