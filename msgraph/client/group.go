package client

const (
	GroupPath    = "groups"
	GroupPathGet = GroupPath + "/%s"
	Unified      = "Unified"
)

type Group struct {
	Id              string   `json:"id,omitempty"`
	Description     string   `json:"description,omitempty"`
	DisplayName     string   `json:"displayName,omitempty"`
	GroupTypes      []string `json:"groupTypes,omitempty"`
	Mail            string   `json:"mail,omitempty"`
	MailEnabled     bool     `json:"mailEnabled"`
	MailNickname    string   `json:"mailNickname,omitempty"`
	SecurityEnabled bool     `json:"securityEnabled"`
	Owners          []string `json:"owners@odata.bind,omitempty"`
}
type GroupCollection struct {
	Groups []Group `json:"value"`
}
