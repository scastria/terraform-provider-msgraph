package client

const (
	UserPath    = "users"
	UserPathGet = UserPath + "/%s"
)

type User struct {
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}
