package auth

import "os"

var AdminKey string

func init() {
	val, ok := os.LookupEnv("ADMIN_KEY")
	if !ok {
		panic("ADMIN_KEY env var must be set")
	}
	AdminKey = val
}

func CheckKey(key string) bool {
	if key == AdminKey {
		return true
	}
	return false
}