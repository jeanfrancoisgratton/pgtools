// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:00
// Original filename: src/cmd/rolesCommands.go

package cmd

import (
	"fmt"
	"os"
	"pgtools/environment"
	"pgtools/logging"
	"pgtools/roles"
	"pgtools/types"

	"github.com/spf13/cobra"
)

var rolesCmd = &cobra.Command{
	Use:     "roles",
	Aliases: []string{"users"},
	Short:   "Manage PostgreSQL roles (users are roles with LOGIN)",
	Long:    "Add, delete, edit PostgreSQL roles. A 'user' is just a role with the LOGIN attribute.",
}

var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List roles",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}
		if nerr := roles.ListRoles(cfg, types.ListMembers, types.ListVerbose); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
	},
}

var roleAddCmd = &cobra.Command{
	Use:   "add <rolename>",
	Short: "Create a new role",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}

		opts := types.RoleOptions{
			Login:       flagTri(types.RoleLogin, types.RoleNoLogin),
			Superuser:   flagTri(types.RoleSuper, types.RoleNoSuper),
			CreateDB:    flagTri(types.RoleCreateDB, types.RoleNoCreateDB),
			CreateRole:  flagTri(types.RoleCreateRole, types.RoleNoCreateRole),
			Inherit:     flagTri(types.RoleInherit, types.RoleNoInherit),
			Replication: flagTri(types.RoleRepl, types.RoleNoRepl),
			BypassRLS:   flagTri(types.RoleBypassRLS, types.RoleNoBypassRLS),
		}
		if types.NewPassword != "" {
			opts.Password = &types.NewPassword
		}

		if nerr := roles.CreateRole(cfg, args[0], opts); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
		fmt.Printf("Role %q created.\n", args[0])
	},
}

var roleDelCmd = &cobra.Command{
	Use:   "delete <rolename>",
	Short: "Drop a role",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}
		if nerr := roles.DropRole(cfg, args[0], types.RoleCascade, types.RoleIfExists); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
		fmt.Printf("Role %q deleted.\n", args[0])
	},
}

var roleEditCmd = &cobra.Command{
	Use:   "edit <rolename>",
	Short: "Alter role attributes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}

		opts := types.RoleOptions{
			Login:       flagTri(types.RoleLogin, types.RoleNoLogin),
			Superuser:   flagTri(types.RoleSuper, types.RoleNoSuper),
			CreateDB:    flagTri(types.RoleCreateDB, types.RoleNoCreateDB),
			CreateRole:  flagTri(types.RoleCreateRole, types.RoleNoCreateRole),
			Inherit:     flagTri(types.RoleInherit, types.RoleNoInherit),
			Replication: flagTri(types.RoleRepl, types.RoleNoRepl),
			BypassRLS:   flagTri(types.RoleBypassRLS, types.RoleNoBypassRLS),
		}
		if types.NewPassword != "" {
			opts.Password = &types.NewPassword
			logging.Infof("Password provided via edit; will apply with ALTER ROLE ... PASSWORD")
		}
		opts.ClearPassword = types.ClearPassword

		if nerr := roles.AlterRole(cfg, args[0], opts); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
		fmt.Printf("Role %q altered.\n", args[0])
	},
}

var rolePassCmd = &cobra.Command{
	Use:   "passwd <rolename>",
	Short: "Change a role's password",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := environment.LoadConfig()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(err.Code)
		}

		var pw *string
		if !types.ClearPassword {
			if types.NewPassword == "" {
				fmt.Println("ERROR: --password is required unless --clear is used")
				os.Exit(300)
			}
			pw = &types.NewPassword
		}
		if nerr := roles.ChangePassword(cfg, args[0], pw, types.ClearPassword); nerr != nil {
			fmt.Printf("%s\n", nerr.Error())
			os.Exit(nerr.Code)
		}
		if types.ClearPassword {
			fmt.Printf("Role %q password cleared.\n", args[0])
		} else {
			fmt.Printf("Role %q password changed.\n", args[0])
		}
	},
}

func init() {
	rolesCmd.AddCommand(rolesListCmd, roleAddCmd, roleDelCmd, roleEditCmd, rolePassCmd)

	// Common attribute flags for add/edit
	for _, c := range []*cobra.Command{roleAddCmd, roleEditCmd} {
		c.Flags().BoolVar(&types.RoleLogin, "login", false, "Set LOGIN")
		c.Flags().BoolVar(&types.RoleNoLogin, "no-login", false, "Set NOLOGIN")
		c.Flags().BoolVar(&types.RoleSuper, "superuser", false, "Set SUPERUSER")
		c.Flags().BoolVar(&types.RoleNoSuper, "no-superuser", false, "Set NOSUPERUSER")
		c.Flags().BoolVar(&types.RoleCreateDB, "createdb", false, "Set CREATEDB")
		c.Flags().BoolVar(&types.RoleNoCreateDB, "no-createdb", false, "Set NOCREATEDB")
		c.Flags().BoolVar(&types.RoleCreateRole, "createrole", false, "Set CREATEROLE")
		c.Flags().BoolVar(&types.RoleNoCreateRole, "no-createrole", false, "Set NOCREATEROLE")
		c.Flags().BoolVar(&types.RoleInherit, "inherit", false, "Set INHERIT")
		c.Flags().BoolVar(&types.RoleNoInherit, "no-inherit", false, "Set NOINHERIT")
		c.Flags().BoolVar(&types.RoleRepl, "replication", false, "Set REPLICATION")
		c.Flags().BoolVar(&types.RoleNoRepl, "no-replication", false, "Set NOREPLICATION")
		c.Flags().BoolVar(&types.RoleBypassRLS, "bypass-rls", false, "Set BYPASSRLS")
		c.Flags().BoolVar(&types.RoleNoBypassRLS, "no-bypass-rls", false, "Set NOBYPASSRLS")
		c.Flags().StringVarP(&types.NewPassword, "password", "p", "", "Set password (omit to leave unchanged)")
		c.Flags().BoolVar(&types.ClearPassword, "clear", false, "Clear password (PASSWORD NULL)")
	}

	// Delete-only flags
	roleDelCmd.Flags().BoolVar(&types.RoleCascade, "cascade", false, "Drop objects owned by the role (CASCADE)")
	roleDelCmd.Flags().BoolVar(&types.RoleIfExists, "if-exists", false, "Do not error if role does not exist")

	// Passwd flags
	rolePassCmd.Flags().StringVarP(&types.NewPassword, "password", "p", "", "New password")
	rolePassCmd.Flags().BoolVar(&types.ClearPassword, "clear", false, "Clear password (PASSWORD NULL)")

	// List flags
	rolesListCmd.Flags().BoolVar(&types.ListMembers, "members", false, "Include role membership mapping")
	rolesListCmd.Flags().BoolVar(&types.ListVerbose, "verbose", false, "Show all boolean attributes")
	rolesListCmd.Flags().BoolVarP(&types.Quiet, "quiet", "q", false, "Quiet output (names only)")
}

// flagTri encodes two opposite flags into a tristate pointer.
// If both (or neither) are set, returns nil (no change).
func flagTri(yes, no bool) *bool {
	if yes && !no {
		t := true
		return &t
	}
	if no && !yes {
		f := false
		return &f
	}
	return nil
}
