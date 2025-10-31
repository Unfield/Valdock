package permissions

import (
	"strings"
)

type Permission string

var (
	InstanceCreate  Permission = "instance:create"
	InstanceStart   Permission = "instance:start"
	InstanceStop    Permission = "instance:stop"
	InstanceRestart Permission = "instance:restart"
	InstanceRead    Permission = "instance:read"
	InstanceWrite   Permission = "instance:write"
	InstanceDelete  Permission = "instance:delete"
	InstanceAdmin   Permission = "instance:admin"

	ACLRead   Permission = "acl:read"
	ACLWrite  Permission = "acl:write"
	ACLDelete Permission = "acl:delete"
	ACLAdmin  Permission = "acl:admin"

	APIKeyCreate Permission = "apikey:read"
	APIKeyRead   Permission = "apikey:read"
	APIKeyWrite  Permission = "apikey:write"
	APIKeyDelete Permission = "apikey:delete"
	APIKeyAdmin  Permission = "apikey:admin"

	RootAdmin Permission = "*:admin"
)

var APIKeyPermissions = []Permission{
	APIKeyCreate,
	APIKeyRead,
	APIKeyWrite,
	APIKeyDelete,
	APIKeyAdmin,
}

var InstancePermissions = []Permission{
	InstanceCreate,
	InstanceStart,
	InstanceStop,
	InstanceRestart,
	InstanceRead,
	InstanceWrite,
	InstanceDelete,
	InstanceAdmin,
}

var ACLPermissions = []Permission{
	ACLRead,
	ACLWrite,
	ACLDelete,
	ACLAdmin,
}

var permissionMap = map[string]Permission{
	string(RootAdmin): RootAdmin,

	string(InstanceCreate):  InstanceCreate,
	string(InstanceStart):   InstanceStart,
	string(InstanceStop):    InstanceStop,
	string(InstanceRestart): InstanceRestart,
	string(InstanceRead):    InstanceRead,
	string(InstanceWrite):   InstanceWrite,
	string(InstanceDelete):  InstanceDelete,
	string(InstanceAdmin):   InstanceAdmin,

	string(ACLRead):   ACLRead,
	string(ACLWrite):  ACLWrite,
	string(ACLDelete): ACLDelete,
	string(ACLAdmin):  ACLAdmin,

	string(APIKeyCreate): APIKeyCreate,
	string(APIKeyRead):   APIKeyRead,
	string(APIKeyWrite):  APIKeyWrite,
	string(APIKeyDelete): APIKeyDelete,
	string(APIKeyAdmin):  APIKeyAdmin,
}

func ParsePermissionString(input string) []Permission {
	permissions := []Permission{}
	splittedString := strings.SplitSeq(strings.ReplaceAll(input, " ", ""), ",")
	for s := range splittedString {
		if perm, ok := permissionMap[s]; ok {
			permissions = append(permissions, perm)
		}
	}
	return permissions
}

func HasOnePermission(has []Permission, want ...Permission) bool {
	if len(want) == 0 {
		return true
	}
	if len(has) == 0 {
		return false
	}

	if has[0] == RootAdmin {
		return true
	}

	lookup := make(map[Permission]struct{})
	for _, p := range has {
		lookup[p] = struct{}{}
	}

	for _, w := range want {
		if _, ok := lookup[w]; ok {
			return true
		}
	}

	return false
}

func HasAllPermissions(has []Permission, want ...Permission) bool {
	if len(want) == 0 {
		return true
	}
	if len(has) == 0 {
		return false
	}

	if has[0] == RootAdmin {
		return true
	}

	lookup := make(map[Permission]struct{}, len(has))
	for _, p := range has {
		lookup[p] = struct{}{}
	}

	for _, w := range want {
		if _, ok := lookup[w]; !ok {
			return false
		}
	}

	return true
}
