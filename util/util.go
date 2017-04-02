package util

import "os"

func ExistsFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
