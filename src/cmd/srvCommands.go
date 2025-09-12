// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 18:01
// Original filename: src/cmd/srvCommands.go

package cmd

import (
	"fmt"
	"os"
	"pgtools/environment"
	"pgtools/srv"
	"pgtools/types"

	hf "github.com/jeanfrancoisgratton/helperFunctions"
	"github.com/spf13/cobra"
)

//var (
//	reloadLocal   bool   // if true: try SQL first, then local SIGHUP fallback
//	reloadPIDFile string // optional local PID file path for fallback
//	reloadDataDir string // optional local data dir for fallback
//)

var srvCmd = &cobra.Command{
	Use:     "srv",
	Aliases: []string{"server"},
	Short:   "server sub-command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Valid subcommands are: { show | backup | restore }")
	},
}

var srvVerCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the database server version",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}

		if srvinfo, err := srv.ShowDBServerVersion(cfg); err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		} else {
			fmt.Printf("Server: %s:%s\nVersion: %s\n",
				hf.Green(srvinfo.ServerName),
				hf.Green(fmt.Sprintf("%v", srvinfo.ServerPort)),
				hf.Green(srvinfo.Version))
		}
	},
}

var srvReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload server configuration",
	Long: `Reload PostgreSQL configuration.

By default this calls SQL: SELECT pg_reload_conf() — works for local or remote servers.
If --local is set, it will attempt a local SIGHUP fallback (using postmaster.pid)
if the SQL call fails. You can pass --pidfile or --datadir to help locate the PID;
if neither is provided, we'll try to discover the data directory via SQL.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}

		if !types.ReloadLocal {
			// Pure SQL path (works for local or remote)
			if nerr := srv.ReloadRemote(cfg); nerr != nil {
				fmt.Printf("%s\n", nerr.Error())
				os.Exit(nerr.Code)
			}
			fmt.Println("Configuration reloaded via SQL.")
			return
		}

		// Local fallback path: try SQL first, then SIGHUP using PID
		local := types.LocalControl{
			DataDir: types.ReloadDataDir,
			PIDFile: types.ReloadPIDFile,
		}

		// If neither PIDFILE nor DATADIR was provided, try to discover DATADIR via SQL.
		if local.PIDFile == "" && local.DataDir == "" {
			if dd, derr := srv.DiscoverDataDirViaSQL(cfg); derr == nil && dd != "" {
				local.DataDir = dd
			}
			// If discovery fails, ReloadLocal will surface a clear error during PID resolution.
		}

		if nerr := srv.ReloadLocal(cfg, local); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
		fmt.Println("Configuration reloaded (SQL with local fallback if needed).")
	},
}

func init() {
	srvCmd.AddCommand(srvVerCmd, srvReloadCmd)

	// Use regular Flags() on the leaf command so they show in help
	srvReloadCmd.Flags().BoolVar(&types.ReloadLocal, "local", false, "Attempt local SIGHUP fallback if SQL reload fails")
	srvReloadCmd.Flags().StringVar(&types.ReloadPIDFile, "pidfile", "", "Path to postmaster.pid for local fallback")
	srvReloadCmd.Flags().StringVar(&types.ReloadDataDir, "datadir", "", "PostgreSQL data directory for local fallback (used to resolve postmaster.pid)")
}
