// pgtool
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/04 16:00
// Original filename: src/cmd/envCommands.go

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"pgtool/environment"
	"pgtool/types"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Valid subcommands are: { list | add | remove | info }")
	},
}

// List config file
var envListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "pgtool env list [directory]",
	Short:   "Lists all env files",
	Run: func(cmd *cobra.Command, args []string) {
		if err := environment.ListEnvironments(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}
	},
}

// Create a config file
var envAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Example: "pgtool env add -e [FILE[.json]]",
	Short:   "Adds the env FILE",
	Long: `The extension (.json) is implied and will be added if missing.
The default defaultEnv.json file will be used if no filename is provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := environment.CreateConfig(); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}
	},
}

// Describe the contents of a config file
var envInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"explain"},
	Example: "pgtool env info FILE1[.json] FILE2[.json]... FILEn[.json]",
	Short:   "Prints the env FILE[12n] information",
	Long:    `You can list as many env files as you wish, here`,
	Run: func(cmd *cobra.Command, args []string) {
		envfiles := []string{types.EnvConfigFile}
		if len(args) != 0 {
			envfiles = args
		}
		if err := environment.ExplainEnvFile(envfiles); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}
	},
}

// Remove a config file
var envRmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove"},
	Example: "pgtool env remove { FILE[.json] | defaultEnv.json }",
	Short:   "Removes the env FILE",
	Run: func(cmd *cobra.Command, args []string) {
		fname := ""
		if len(args) == 0 {
			fname = types.EnvConfigFile
		} else {
			fname = args[0]
		}
		if err := environment.RemoveEnvFile(fname); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}
	},
}

func init() {
	envCmd.AddCommand(envRmCmd, envInfoCmd, envAddCmd, envListCmd)

}
