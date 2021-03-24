package client

import "strings"

const (
	GroupMemberPath       = "groups/%s/members"
	GroupMemberPathCreate = GroupMemberPath + "/$ref"
	GroupMemberPathDelete = GroupMemberPath + "/%s/$ref"
)

type GroupMember struct {
	GroupId  string `json:"-"`
	MemberId string `json:"-"`
	OdataId  string `json:"@odata.id,omitempty"`
}

func (grpmem *GroupMember) GroupMemberEncodeId() string {
	return grpmem.GroupId + IdSeparator + grpmem.MemberId
}

func GroupMemberDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
