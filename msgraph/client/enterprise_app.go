package client

const (
	EnterpriseAppPath    = "servicePrincipals"
	EnterpriseAppPathGet = EnterpriseAppPath + "/%s"
)

type EnterpriseApp struct {
	Id             string `json:"id,omitempty"`
	AppId          string `json:"appId,omitempty"`
	DisplayName    string `json:"displayName,omitempty"`
	AppDisplayName string `json:"appDisplayName,omitempty"`
	LoginUrl       string `json:"loginUrl,omitempty"`
	LogoutUrl      string `json:"logoutUrl,omitempty"`
}
type EnterpriseAppCollection struct {
	EnterpriseApps []EnterpriseApp `json:"value"`
}
