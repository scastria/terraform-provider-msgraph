package client

const (
	UserPath    = "users"
	UserPathGet = UserPath + "/%s"
)

type User struct {
	DisplayName       string `json:"displayName,omitempty"`
	Mail              string `json:"mail,omitempty"`
	UserPrincipalName string `json:"userPrincipalName,omitempty"`
}
