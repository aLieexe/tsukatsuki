package utils

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// func CheckServerReachable() error

func GetIPVersion(ip net.IP) string {
	if ip.To4() != nil {
		return "v4"
	} else {
		return "v6"
	}
}

// expect if public dns can actually reach it
func IsDomainConfigured(domainName string) bool {
	_, err := net.LookupHost(domainName)
	return err == nil
}

func IpValidator(input string) error {
	if net.ParseIP(input) == nil {
		return fmt.Errorf("must be a valid IP address")
	}
	return nil
}

func PortValidator(input string) error {
	parsed, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("must be an integer")
	}

	if parsed < 1 || parsed > 65535 {
		return fmt.Errorf("%v is not a valid port number (1 - 65535)", parsed)
	}

	return nil
}

func SiteAddressValidator(input string) error {
	_, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("must be a valid address")
	}

	return nil
}

func GetMainFileLocation() string {
	cmd := exec.Command("sh", "-c", `find . -type f -name "*.go" -exec grep -m1 -H '^func main()' {} + | grep -v '^[[:space:]]*//' | head -n1 | cut -d: -f1`)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

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
