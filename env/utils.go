package env

import "regexp"

var podNameSuffix = regexp.MustCompile(`^(.*)(\-[a-z0-9]{10})(\-[a-z0-9]{5})$`)

func HostNameToServiceName(hostName string) string {
	return podNameSuffix.ReplaceAllString(hostName, "$1")
}
