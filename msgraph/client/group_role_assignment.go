package client

import "strings"

const (
	GroupRoleAssignmentPath    = "groups/%s/appRoleAssignments"
	GroupRoleAssignmentPathGet = GroupRoleAssignmentPath + "/%s"
)

type GroupRoleAssignment struct {
	GroupId     string `json:"-"`
	Id          string `json:"id,omitempty"`
	ResourceId  string `json:"resourceId,omitempty"`
	PrincipalId string `json:"principalId,omitempty"`
	AppRoleId   string `json:"appRoleId,omitempty"`
}

func (gra *GroupRoleAssignment) GroupRoleAssignmentEncodeId() string {
	return gra.GroupId + IdSeparator + gra.Id
}

func GroupRoleAssignmentDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
