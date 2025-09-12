// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 06:59
// Original filename: src/cmd/listCommands.go

package cmd

import (
	"fmt"
	"os"
	"pgtools/environment"
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

var listDBCmd = &cobra.Command{
	Use:     "dbs",
	Aliases: []string{"databases"},
	Short:   "List all accessible databases",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		if _, nerr := show.ListDatabases(cfg, types.SortBySize); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
	},
}

func init() {
	showCmd.AddCommand(listDBCmd)

	listDBCmd.Flags().BoolVarP(&types.SortBySize, "sort-size", "s", false, "Sort databases by size (largest first)")
}
