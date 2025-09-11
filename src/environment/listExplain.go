// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/05 16:49
// Original filename: src/environment/listExplain.go

package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"pgtools/types"
	"strconv"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v2"
	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ListEnvironments() *ce.CustomError {
	var err error
	var dirFH *os.File
	var finfo, fileInfos []os.FileInfo

	if dirFH, err = os.Open(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtools")); err != nil {
		return &ce.CustomError{Code: 15, Title: "Unable to read config directory", Message: err.Error()}
	}

	if fileInfos, err = dirFH.Readdir(0); err != nil {
		return &ce.CustomError{Code: 16, Title: "Unable to read files in config directory", Message: err.Error()}
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") && !strings.HasPrefix(info.Name(), "sample") {
			finfo = append(finfo, info)
		}
	}

	if err != nil {
		return &ce.CustomError{Code: 99, Title: "Undefined error", Message: err.Error()}
	}

	fmt.Printf("Number of environment files: %s\n", hf.Green(fmt.Sprintf("%d", len(finfo))))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Environment file", "File size", "Modification time"})

	for _, fi := range finfo {
		t.AppendRow([]interface{}{hf.Green(fi.Name()), hf.Green(hf.SI(uint64(fi.Size()))),
			hf.Green(fmt.Sprintf("%v", fi.ModTime().Format("2006/01/02 15:04:05")))})
	}
	t.SortBy([]table.SortBy{
		{Name: "Environment file", Mode: table.Asc},
		{Name: "File size", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}

func ExplainEnvFile(envfiles []string) *ce.CustomError {
	oldEnvFile := types.EnvConfigFile

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Environment file", "DB server host", "DB server port", "DB user", "DB password",
		"SSL enabled", "SSL cert", "SSL key", "Description"})

	for _, envfile := range envfiles {
		if !strings.HasSuffix(envfile, ".json") {
			envfile += ".json"
		}
		types.EnvConfigFile = envfile

		// error code 99 is specifically taylored for this use; it means that defaultEnv does not exist, and that
		// we can continue. otherwise we abort
		if e, err := LoadConfig(); err != nil && err.Code != 99 {
			types.EnvConfigFile = oldEnvFile
			return err
		} else {
			sslcert := filepath.Base(e.SSLCert)
			sslkey := filepath.Base(e.SSLKey)
			if sslcert == "." {
				sslcert = "n/a"
			}
			if sslkey == "." {
				sslkey = "n/a"
			}
			t.AppendRow([]interface{}{hf.Green(envfile), hf.Green(e.Host), hf.Green(strconv.Itoa(e.Port)),
				hf.Green(e.User), hf.Yellow("*ENCODED*"),
				hf.Green(e.SSLMode), hf.Green(sslcert), hf.Green(sslkey), hf.Green(e.Description)})
		}
	}
	t.SortBy([]table.SortBy{
		{Name: "Environment file", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	types.EnvConfigFile = oldEnvFile
	return nil
}
