package client

import "strings"

const (
	EnterpriseAppRolePath              = "servicePrincipals/%s/appRoles"
	EnterpriseAppRoleDefaultAccessId   = "00000000-0000-0000-0000-000000000000"
	EnterpriseAppRoleDefaultAccessName = "Default Access"
)

type EnterpriseAppRole struct {
	EnterpriseAppId string `json:"-"`
	Id              string `json:"id,omitempty"`
	DisplayName     string `json:"displayName,omitempty"`
	Description     string `json:"description,omitempty"`
}
type EnterpriseAppRoleCollection struct {
	EnterpriseAppRoles []EnterpriseAppRole `json:"value"`
}

func (ear *EnterpriseAppRole) EnterpriseAppRoleEncodeId() string {
	return ear.EnterpriseAppId + IdSeparator + ear.Id
}

func EnterpriseAppRoleDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
