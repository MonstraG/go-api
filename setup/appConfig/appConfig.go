package appConfig

import (
	"encoding/json"
	"flag"
	"go-server/setup/myLog"
	"os"
)

type AppConfig struct {
	Host             string      `json:"host"`
	DatabaseFile     string      `json:"databaseFile"`
	DefaultUser      DefaultUser `json:"defaultUser"`
	GuestUser        DefaultUser `json:"guestUser"`
	VpsLoginLink     string      `json:"vpsLoginLink"`
	JWTSecret        string      `json:"jwtSecret"`
	SongsFolder      string      `json:"songsFolder"`
	WebSocketOrigins []string    `json:"webSocketOrigins"`
}

type DefaultUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ReadConfig() AppConfig {
	configPath := readConfigPath()

	configJson, err := os.ReadFile(configPath)
	if err != nil {
		myLog.Fatal.Logf("Failed to read config file in '%s':\n%v\n", configPath, err)
	}

	var config AppConfig
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		myLog.Fatal.Logf("Failed to unmarshall config file in '%s':\n%v\n", configPath, err)
	}

	return config
}

func readConfigPath() string {
	configPathVar := flag.String("config", "config.json", "Path to json config for the server")
	flag.Parse()
	myLog.Info.Logf("Loading config from \"%s\"", *configPathVar)
	return *configPathVar
}
