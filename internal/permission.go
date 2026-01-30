package internal

import (
	"strings"
)

func permissionMatch(pattern string, permission string) bool {
	// Fast path
	if pattern == "*" {
		return true
	}

	// Enforce at most one wildcard
	if strings.Count(pattern, "*") > 1 {
		return false
	}

	// No wildcard → exact match
	if !strings.Contains(pattern, "*") {
		return pattern == permission
	}

	// Exactly one wildcard
	prefix, suffix, _ := strings.Cut(pattern, "*")
	return strings.HasPrefix(permission, prefix) &&
		strings.HasSuffix(permission, suffix)
}

func HasPermission(info AuthInfo, requiredPerm string) bool {
	for _, perm := range info.Permissions {
		if permissionMatch(perm, requiredPerm) {
			return true
		}
	}

	return false
}
