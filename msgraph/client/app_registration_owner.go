package client

import "strings"

const (
	AppRegistrationOwnerPath       = "applications/%s/owners"
	AppRegistrationOwnerPathCreate = AppRegistrationOwnerPath + "/$ref"
	AppRegistrationOwnerPathDelete = AppRegistrationOwnerPath + "/%s/$ref"
)

type AppRegistrationOwner struct {
	AppRegistrationId string `json:"-"`
	OwnerId           string `json:"-"`
	OdataId           string `json:"@odata.id,omitempty"`
}

func (arown *AppRegistrationOwner) AppRegistrationOwnerEncodeId() string {
	return arown.AppRegistrationId + IdSeparator + arown.OwnerId
}

func AppRegistrationOwnerDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
