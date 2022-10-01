package permissions

import "github.com/samber/lo"

type Claims = map[string]string
type Permissions map[string][]Claims

func (perms Permissions) GetPermissions(providedClaims Claims) []string {
	res := []string{}
	for permission, requiredClaimsSlice := range perms {
		if lo.SomeBy(requiredClaimsSlice, func(requiredClaims Claims) bool {
			for k, v := range requiredClaims {
				if providedClaims[k] != v {
					return false
				}
			}
			return true
		}) {
			res = append(res, permission)
		}
	}
	return res
}
