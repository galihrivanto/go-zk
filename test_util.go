package gozk

import "os"

func getEnvar(key string, vars ...string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		if len(vars) > 0 {
			return vars[0]
		}
	}

	return v
}
