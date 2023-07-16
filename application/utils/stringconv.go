package utils

func TrimStr(v string) string {
	if len(v) > 20 {
		return v[0:20]
	}
	return v
}
