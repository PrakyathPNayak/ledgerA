package service

import "strings"

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(strings.ToLower(err.Error()), "duplicate key")
}
