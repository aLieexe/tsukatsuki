package utils

import (
	"os"
	"path/filepath"
)

// func CheckServerReachable() error

// func ValidateIP() error
// func GetIPVersion() error

// // should be between debian, RHEL? Idk if there is difference between CentOS or Alma or other
// func GetDistribution() error

// func ValidatePort() error

// func ReplaceWithHyphens ()

// func CheckDomainDNSConfiguration ()

func GetProjectDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(dir)
}
