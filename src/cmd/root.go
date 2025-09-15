// pgtools
// src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"pgtools/types"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pgtools",
	Short:   "PostgreSQL client utilities",
	Version: "1.70.00 (2025.09.15)",
}

// Shows changelog
var clCmd = &cobra.Command{
	Use:     "changelog",
	Aliases: []string{"cl"},
	Short:   "Shows the Changelog",
	Run: func(cmd *cobra.Command, args []string) {
		changeLog()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.DisableAutoGenTag = true
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.AddCommand(completionCmd, clCmd, envCmd, dbCmd, rolesCmd, srvCmd, showCmd, confCmd)

	rootCmd.PersistentFlags().StringVarP(&types.LogLevel, "loglevel", "l", "none", "Log level: none|debug|info|error")
	rootCmd.PersistentFlags().BoolVarP(&types.DebugMode, "debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&types.EnvConfigFile, "env", "e", "defaultEnv.json", "Default environment configuration file; this is a per-user setting.")
}

func changeLog() {
	//fmt.Printf("\x1b[2J")
	fmt.Printf("\x1bc")

	fmt.Println("CHANGELOG")
	fmt.Println("=========")
	fmt.Println()

	fmt.Print(`
VERSION			DATE			COMMENT
-------			----			-------
1.70.00			2025.09.15		Completed the conf subcommand
1.60.00			2025.09.15		Completed the show subcommand
1.50.00			2025.09.12		Completed the db create and db drop subcommands
1.40.00			2025.09.11		Completed the srv subcommand
1.30.00			2025.09.11		Completed the roles subcommand
1.21.10			2025.09.07		Updated to GO 1.25.1
1.21.00			2025.07.24		Fixed a nil-pointer issue with db backup -u
1.20.00			2025.07.24		Error handling is firmed-up, more consistent error codes (Phase 1)
1.10.00			2025.07.15		Fixed double-quote issue, completed constraints management
1.07.00			2025.07.14		pgtool db backup is now near-parity with pg_dumpall
1.05.00			2025.07.09		Happy birthday, Liliane. Testing the restore subcommand
1.00.00			2025.07.04		Initial release, backup subcommand tested
`)
}
