// pgtool
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/05 18:57
// Original filename: src/cmd/pgCommands.go

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"pgtool/db"
	"pgtool/environment"
	"pgtool/types"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Environment sub-command",
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
			fmt.Println(os.Stderr, "Failed to load config:", err)
			os.Exit(err.Code)
		}
		if err := db.BackupDatabase(cfg, args); err != nil {
			fmt.Println(os.Stderr, "Failed to backup database:", err)
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
			fmt.Println("pgtool restore ARCHIVE_NAME")
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

var dbVerCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the database server version",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}

		if ver, err := db.ShowDBServerVersion(cfg); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		} else {
			fmt.Printf("Server version: %s\n", ver)
		}
	},
}

func init() {
	dbCmd.AddCommand(listCmd, backupCmd, restoreCmd, dbVerCmd)
	rootCmd.AddCommand(listCmd, backupCmd, restoreCmd)

	rootCmd.PersistentFlags().StringVarP(&types.LogLevel, "loglevel", "l", "none", "Log level: none|debug|info|error")
	backupCmd.PersistentFlags().BoolVarP(&types.UserRoles, "users", "u", false, "Backup global users/roles only")
	backupCmd.PersistentFlags().BoolVarP(&types.AllDBs, "all", "a", false, "Backup all databases")
	backupCmd.MarkFlagsMutuallyExclusive("all", "users")

	restoreCmd.PersistentFlags().StringVarP(&types.LogLevel, "loglevel", "l", "error", "Log level: debug|info|error")
	restoreCmd.PersistentFlags().BoolVarP(&types.UserRoles, "users", "u", false, "Backup global users/roles only")

	listCmd.PersistentFlags().BoolVarP(&types.Quiet, "quiet", "q", false, "Silent output")

}
