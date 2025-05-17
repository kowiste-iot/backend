package domain

import (
	"fmt"
	"strings"
)

const (
	defaultName string = "-policy"
)
const (
	TypeRole  string = "role"
	Enforcing string = "ENFORCING"
)

const (
	TypeScope    string = "scope"
	TypeResource string = "resource"
)
const (
	DecisionUnanimous   string = "UNANIMOUS"
	DecisionAffirmative string = "AFFIRMATIVE"
)
const (
	LogicPositive string = "POSITIVE"
)
func PolicyName(roleName string) string {
	return fmt.Sprintf("%s"+defaultName, roleName)
}
func PolicyToRole(policyName string) string {
	return strings.Replace(policyName, defaultName, "", -1)
}
