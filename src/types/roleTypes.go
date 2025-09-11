// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/09/11 13:35
// Original filename: src/types/roleTypes.go

package types

var (
	RoleLogin        bool
	RoleNoLogin      bool
	RoleSuper        bool
	RoleNoSuper      bool
	RoleCreateDB     bool
	RoleNoCreateDB   bool
	RoleCreateRole   bool
	RoleNoCreateRole bool
	RoleInherit      bool
	RoleNoInherit    bool
	RoleRepl         bool
	RoleNoRepl       bool
	RoleBypassRLS    bool
	RoleNoBypassRLS  bool

	RoleCascade  bool
	RoleIfExists bool

	NewPassword   string
	ClearPassword bool

	// list flags
	ListMembers bool
	ListVerbose bool
)

type RoleOptions struct {
	Login       *bool
	Superuser   *bool
	CreateDB    *bool
	CreateRole  *bool
	Inherit     *bool
	Replication *bool
	BypassRLS   *bool

	// Password handling:
	//   Password == nil        -> no change (ALTER) / omit (CREATE)
	//   ClearPassword == true  -> ALTER ROLE ... PASSWORD NULL
	//   Password != nil        -> set to provided value
	Password      *string
	ClearPassword bool
}
