// pgtool
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/04 12:34
// Original filename: src/environment/environment.go

package environment

import (
	"encoding/json"
	"fmt"
	ce "github.com/jeanfrancoisgratton/customError/v2"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"os"
	"path/filepath"
	"pgtool/types"
	"strings"
)

// Loads the configuration file if -e is not passed, it defaults to $HOME/.environment/JFG/pgtool/defaultEnv.json
func LoadConfig() (*types.DBConfig, *ce.CustomError) {
	if !strings.HasSuffix(types.EnvConfigFile, ".json") {
		types.EnvConfigFile += ".json"
	}
	path := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtool", types.EnvConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
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
	dbc.Host = hf.GetStringValFromPrompt("PGSQL server hostname: ")
	dbc.Port = hf.GetIntValFromPrompt("PGSQL server port: ")
	dbc.User = hf.GetStringValFromPrompt("PGSQL server username: ")
	dbc.Password = hf.EncodeString(hf.GetPassword("Please enter the user's password: ", types.DebugMode), "")
	sslmode := hf.GetBoolValFromPrompt("PGSQL SSL mode (t/f): ")
	if sslmode {
		dbc.SSLMode = "require"
	} else {
		dbc.SSLMode = "disable"
	}
	dbc.SSLCert = hf.GetStringValFromPrompt("[optional] PGSQL SSL certificate: ")
	dbc.SSLKey = hf.GetStringValFromPrompt("[optional] PGSQL SSL key: ")

	jStream, err := json.MarshalIndent(dbc, "", "  ")
	if err != nil {
		return &ce.CustomError{Code: 11, Title: "Error marshalling JSON", Message: err.Error()}
	}
	types.EnvConfigFile = filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtool", types.EnvConfigFile)

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
	if err := os.Remove(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtool", envfile)); err != nil {
		return &ce.CustomError{Code: 13, Title: "Error removing " + envfile, Message: err.Error()}
	}

	fmt.Printf("%s removed succesfully\n", envfile)
	return nil
}
