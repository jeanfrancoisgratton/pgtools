// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/04 12:34
// Original filename: src/environment/environment.go

package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"pgtools/types"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v2"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
)

// Loads the configuration file if -e is not passed, it defaults to $HOME/.config/JFG/pgtools/defaultEnv.json
func LoadConfig() (*types.DBConfig, *ce.CustomError) {
	if !strings.HasSuffix(types.EnvConfigFile, ".json") {
		types.EnvConfigFile += ".json"
	}
	path := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtools", types.EnvConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		_, a := os.Stat(path)
		if a != nil && !os.IsNotExist(a) && types.EnvConfigFile == "defaultEnv.json" {
			// The code is ignored everywhere except in the for__loop in pgtools env info
			return nil, &ce.CustomError{Title: "Failed to read environment file", Message: err.Error(), Code: 99}
		}
		return nil, &ce.CustomError{Title: "Failed to read environment file", Message: err.Error()}
	}

	var cfg types.DBConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, &ce.CustomError{Code: 10, Title: "Failed to marshal JSON", Message: err.Error()}
	}
	cfg.Password = hf.DecodeString(cfg.Password, "")
	return &cfg, nil
}

// Prompts the user to create an environment environment file, and saves it
func CreateConfig() *ce.CustomError {
	var dbc types.DBConfig

	dbc.Description = hf.GetStringValFromPrompt("[optional] Brief description or comment: ")
	dbc.Host = hf.GetStringValFromPrompt("PGSQL server hostname: ")
	dbc.Port = hf.GetIntValFromPrompt("PGSQL server port: ")
	dbc.User = hf.GetStringValFromPrompt("PGSQL server username: ")
	dbc.Password = hf.EncodeString(hf.GetPassword("Please enter the user's password: ", types.DebugMode), "")
	sslmode := hf.GetBoolValFromPrompt("PGSQL SSL mode (t/f): ")
	if sslmode {
		dbc.SSLMode = "require"
		dbc.SSLCert = hf.GetStringValFromPrompt("[optional] Path to the PGSQL SSL certificate: ")
		dbc.SSLKey = hf.GetStringValFromPrompt("[optional] Path to the PGSQL SSL key: ")
	} else {
		dbc.SSLMode = "disable"
	}
	if sslmode && (dbc.SSLCert == "" || dbc.SSLKey == "") {
		return &ce.CustomError{Code: 12, Title: "SSL certificate and key are required"}
	}

	jStream, err := json.MarshalIndent(dbc, "", "  ")
	if err != nil {
		return &ce.CustomError{Code: 11, Title: "Error marshalling JSON", Message: err.Error()}
	}
	types.EnvConfigFile = filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtools", types.EnvConfigFile)

	if err = os.WriteFile(types.EnvConfigFile, jStream, 0700); err != nil {
		return &ce.CustomError{Code: 12,
			Title:   fmt.Sprintf("Error writing the config file %s", types.EnvConfigFile),
			Message: err.Error()}
	}
	return nil
}

// Remove the config file
func RemoveEnvFile(envfile string) *ce.CustomError {
	if !strings.HasSuffix(envfile, ".json") {
		envfile += ".json"
	}
	if err := os.Remove(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtools", envfile)); err != nil {
		return &ce.CustomError{Code: 13, Title: "Error removing " + envfile, Message: err.Error()}
	}

	fmt.Printf("%s removed succesfully\n", envfile)
	return nil
}
