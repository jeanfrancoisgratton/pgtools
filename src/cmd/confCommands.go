// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/12 06:59
// Original filename: src/cmd/showCommands.go

package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"pgtools/conf"
	"pgtools/shared"

	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Configuration (SHOW ALL / ALTER SYSTEM) helpers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Valid subcommands are: { list | get | set }")
	},
}

var confListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all PostgreSQL configuration parameters (SHOW ALL)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, cerr := shared.GetPool(ctx)
		if cerr != nil {
			fmt.Println(cerr.Error())
			os.Exit(1)
		}
		defer pool.Close()

		rows, err := conf.CollectAll(ctx, pool)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		conf.Render(rows)
	},
}

var confGetCmd = &cobra.Command{
	Use:   "get <key1> [key2 ...]",
	Short: "Get one or more configuration parameters",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, cerr := shared.GetPool(ctx)
		if cerr != nil {
			fmt.Println(cerr.Error())
			os.Exit(1)
		}
		defer pool.Close()

		rows, err := conf.CollectByNames(ctx, pool, args)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if len(rows) == 0 {
			fmt.Println("No matching settings.")
			return
		}
		conf.Render(rows)
	},
}

var confSetCmd = &cobra.Command{
	Use:   "set key = value",
	Short: "Set a configuration parameter using ALTER SYSTEM, then reload",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Accept: key "=" value | key=value | key value
		key := ""
		val := ""

		if len(args) >= 3 && args[1] == "=" {
			key, val = args[0], strings.Join(args[2:], " ")
		} else if len(args) >= 1 {
			joined := strings.Join(args, " ")
			i := strings.Index(joined, "=")
			if i > 0 {
				key = strings.TrimSpace(joined[:i])
				val = strings.TrimSpace(joined[i+1:])
			} else if len(args) == 2 {
				key, val = args[0], args[1]
			}
		}

		// basic name validation (Postgres GUCs allow letters/digits/underscore/dot)
		if key == "" || val == "" || !regexp.MustCompile(`^[a-zA-Z0-9_.]+$`).MatchString(key) {
			fmt.Println("Usage: pgtools conf set key = value")
			os.Exit(2)
		}

		ctx, cancel := shared.CancellableContext()
		defer cancel()

		pool, cerr := shared.GetPool(ctx)
		if cerr != nil {
			fmt.Println(cerr.Error())
			os.Exit(1)
		}
		defer pool.Close()

		if err := conf.Set(ctx, pool, key, val); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Set %s = %s and reloaded configuration.\n", key, val)
	},
}

func init() {
	confCmd.AddCommand(confListCmd, confGetCmd, confSetCmd)
	rootCmd.AddCommand(confCmd)
}
