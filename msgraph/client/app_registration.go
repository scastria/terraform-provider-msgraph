package client

const (
	AppRegistrationPath    = "applications"
	AppRegistrationPathGet = AppRegistrationPath + "/%s"
)

type AppRegistration struct {
	Id          string `json:"id,omitempty"`
	AppId       string `json:"appId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
type AppRegistrationCollection struct {
	AppRegistrations []AppRegistration `json:"value"`
}
