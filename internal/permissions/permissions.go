package permissions

import (
	"encoding/json"
	"os"

	"github.com/samber/lo"
)

type Claims = map[string]string
type Permissions map[string][]Claims

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
