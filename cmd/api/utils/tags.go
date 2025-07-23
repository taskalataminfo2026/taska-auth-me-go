package utils

func Merge(tags1 []string, tags2 ...string) []string {
	return append(tags1, tags2...)
}
