package client

import "strings"

const (
	GroupOwnerPath       = "groups/%s/owners"
	GroupOwnerPathCreate = GroupOwnerPath + "/$ref"
	GroupOwnerPathDelete = GroupOwnerPath + "/%s/$ref"
)

type GroupOwner struct {
	GroupId string `json:"-"`
	OwnerId string `json:"-"`
	OdataId string `json:"@odata.id,omitempty"`
}

func (grpown *GroupOwner) GroupOwnerEncodeId() string {
	return grpown.GroupId + IdSeparator + grpown.OwnerId
}

func GroupOwnerDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
