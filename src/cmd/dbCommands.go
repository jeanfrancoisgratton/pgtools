// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/05 18:57
// Original filename: src/cmd/dbCommands.go

package cmd

import (
	"fmt"
	"os"
	"pgtools/db"
	"pgtools/environment"
	"pgtools/types"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:     "db",
	Aliases: []string{"database"},
	Short:   "Database sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Valid subcommands are: { list | backup | restore }")
	},
}
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all accessible databases",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		if _, nerr := db.ListDatabases(cfg); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
	},
}

var backupCmd = &cobra.Command{
	Use:     "backup [-a] db1 [db2 ...] output.tar[.gz]",
	Short:   "Backup one or more databases to a tarball archive",
	Aliases: []string{"dump"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Println("Failed to load config:", err.Error())
			os.Exit(err.Code)
		}
		if err := db.BackupDatabase(cfg, args); err != nil {
			fmt.Println("Failed to backup database:", err.Error())
			os.Exit(err.Code)
		}
	},
}

var restoreCmd = &cobra.Command{
	Use:     "restore",
	Short:   "Restores one or more databases",
	Aliases: []string{"load"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("pgtools restore ARCHIVE_NAME")
			os.Exit(1)
		}
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}

		if err := db.RestoreDatabase(cfg, args); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}
	},
}

func init() {
	dbCmd.AddCommand(listCmd, backupCmd, restoreCmd)

	backupCmd.PersistentFlags().BoolVarP(&types.UserRoles, "users", "u", false, "Backup global users/roles only")
	backupCmd.PersistentFlags().BoolVarP(&types.AllDBs, "all", "a", false, "Backup all databases")
	backupCmd.MarkFlagsMutuallyExclusive("all", "users")

	restoreCmd.PersistentFlags().StringVarP(&types.LogLevel, "loglevel", "l", "error", "Log level: debug|info|error")
	restoreCmd.PersistentFlags().BoolVarP(&types.UserRoles, "users", "u", false, "Backup global users/roles only")

	listCmd.PersistentFlags().BoolVarP(&types.Quiet, "quiet", "q", false, "Silent output")
}
