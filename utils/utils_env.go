package utils

import "os"

//GetEnv get env variable by key, return defaultValue if empty
func GetEnv(key string, defaultValue string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return defaultValue
	}
	return v
}
