// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 06:59
// Original filename: src/cmd/showCommands.go

package cmd

import (
	"fmt"
	"os"
	"pgtools/environment"
	"pgtools/shared"
	"pgtools/show"
	"pgtools/types"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Valid subcommands are: { dbs | sessions | schemas | stats }")
	},
}

var showDBsCmd = &cobra.Command{
	Use:     "dbs",
	Aliases: []string{"databases"},
	Short:   "List all accessible databases",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		if _, nerr := show.ShowDatabases(cfg, types.SortBySize); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
	},
}

var showSchemasCmd = &cobra.Command{
	Use:   "schemas",
	Short: "List all schemas",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, err := shared.GetPool(ctx)
		if err != nil {
			return err
		}
		defer pool.Close()

		return show.ShowSchemas(ctx, pool)
	},
}

var showTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "List all tables with sizes and row estimates",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, err := shared.GetPool(ctx)
		if err != nil {
			return err
		}
		defer pool.Close()

		return show.ShowTables(ctx, pool)
	},
}

var showSessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "List active sessions (pg_stat_activity)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, err := shared.GetPool(ctx)
		if err != nil {
			return err
		}
		defer pool.Close()

		return show.ShowSessions(ctx, pool)
	},
}

func init() {
	showCmd.AddCommand(showDBsCmd, showSchemasCmd, showTablesCmd, showSessionsCmd)
	showDBsCmd.Flags().BoolVarP(&types.SortBySize, "sort-size", "s", false, "Sort databases by size (largest first)")
}
