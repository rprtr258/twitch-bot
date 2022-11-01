package permissions

import (
	"encoding/json"
	"os"

	"github.com/samber/lo"
)

type Claims = map[string]string
type Permissions map[string][]Claims
type PermissionsList map[string]struct{}

func LoadFromJSONFile(filename string) (Permissions, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data Permissions
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (perms Permissions) GetPermissions(providedClaims Claims) PermissionsList {
	res := make(PermissionsList)
	for permission, requiredClaimsSlice := range perms {
		if lo.SomeBy(requiredClaimsSlice, func(requiredClaims Claims) bool {
			for k, v := range requiredClaims {
				if providedClaims[k] != v {
					return false
				}
			}
			return true
		}) {
			res[permission] = struct{}{}
		}
	}
	return res
}

func (perms PermissionsList) Has(permissions ...string) bool {
	return lo.EveryBy(permissions, func(permission string) bool {
		_, ok := perms[permission]
		return ok
	})
}
