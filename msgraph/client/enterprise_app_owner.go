package client

import "strings"

const (
	EnterpriseAppOwnerPath       = "servicePrincipals/%s/owners"
	EnterpriseAppOwnerPathCreate = EnterpriseAppOwnerPath + "/$ref"
	EnterpriseAppOwnerPathDelete = EnterpriseAppOwnerPath + "/%s/$ref"
)

type EnterpriseAppOwner struct {
	EnterpriseAppId string `json:"-"`
	OwnerId         string `json:"-"`
	OdataId         string `json:"@odata.id,omitempty"`
}

func (eaown *EnterpriseAppOwner) EnterpriseAppOwnerEncodeId() string {
	return eaown.EnterpriseAppId + IdSeparator + eaown.OwnerId
}

func EnterpriseAppOwnerDecodeId(s string) (string, string) {
	tokens := strings.Split(s, IdSeparator)
	return tokens[0], tokens[1]
}
