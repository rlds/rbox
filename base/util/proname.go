package util

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetProName() string {
	_, name := filepath.Split(os.Args[0])
	return name
}

func GetUserInfo() (name, hpath string) {
	user, err := user.Current()
	if err != nil {
		return
	}
	return user.Name, user.HomeDir
}
