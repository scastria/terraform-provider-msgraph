package client

const (
	EnterpriseAppPath    = "servicePrincipals"
	EnterpriseAppPathGet = EnterpriseAppPath + "/%s"
	IntegratedApp        = "WindowsAzureActiveDirectoryIntegratedApp"
)

type EnterpriseApp struct {
	Id          string    `json:"id,omitempty"`
	AppId       string    `json:"appId,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	LoginUrl    string    `json:"loginUrl,omitempty"`
	LogoutUrl   string    `json:"logoutUrl,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	AppRoles    []AppRole `json:"appRoles,omitempty"`
}
type AppRoles struct {
	AppRoles []AppRole `json:"appRoles,omitempty"`
}
type AppRole struct {
	Id                 string   `json:"id,omitempty"`
	AllowedMemberTypes []string `json:"allowedMemberTypes,omitempty"`
	DisplayName        string   `json:"displayName,omitempty"`
	Description        string   `json:"description,omitempty"`
	IsEnabled          bool     `json:"isEnabled"`
}
type EnterpriseAppCollection struct {
	EnterpriseApps []EnterpriseApp `json:"value"`
}
