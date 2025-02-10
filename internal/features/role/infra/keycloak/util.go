package keycloak

import (
	"fmt"
	"strings"
)
const (
	defaultName string = "-policy"
)
func policyName(roleName string) string {
	return fmt.Sprintf("%s"+defaultName, roleName)
}
func policyToRole(policyName string) string {
	return strings.Replace(policyName, defaultName, "", -1)
}
