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
	Aliases: []string{"databases", "db"},
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
	Use:   "schemas [DB ...]",
	Short: "List schemas across one or more databases as a single table",
	Long:  "List all non-system schemas for the specified database(s) in a single aggregated table. If no DBs are specified, defaults to 'postgres'.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		dbs := args
		if len(dbs) == 0 {
			dbs = []string{"postgres"}
		}

		var all []types.SchemaRow

		for _, dbname := range dbs {
			pool, perr := shared.GetPoolForDB(ctx, dbname)
			if perr != nil {
				fmt.Println(perr.Error())
				os.Exit(1)
			}

			rows, cerr := show.CollectSchemas(ctx, pool, dbname)
			pool.Close()
			if cerr != nil {
				fmt.Println(cerr.Error())
				os.Exit(1)
			}

			all = append(all, rows...)
		}

		show.RenderSchemas(all)
	},
}

var showTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Show all tables with sizes and row estimates",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, err := shared.GetPool(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer pool.Close()

		if err := show.ShowTables(ctx, pool); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var showSessionsCmd = &cobra.Command{
	Use:     "sessions",
	Aliases: []string{"activity"},
	Short:   "Show active sessions (pg_stat_activity)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, err := shared.GetPool(ctx)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer pool.Close()

		if sserr := show.ShowSessions(ctx, pool); sserr != nil {
			fmt.Println(sserr.Error())
			os.Exit(1)
		}
	},
}

func init() {
	showCmd.AddCommand(showDBsCmd, showSchemasCmd, showTablesCmd, showSessionsCmd)
	showDBsCmd.Flags().BoolVarP(&types.SortBySize, "sort-size", "s", false, "Sort databases by size (largest first)")
}
