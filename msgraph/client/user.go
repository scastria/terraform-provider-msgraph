package client

const (
	UserPath    = "users"
	UserPathGet = UserPath + "/%s"
)

type User struct {
	Id                string `json:"id,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	GivenName         string `json:"givenName,omitempty"`
	Surname           string `json:"surname,omitempty"`
	JobTitle          string `json:"jobTitle,omitempty"`
	Mail              string `json:"mail,omitempty"`
	UserPrincipalName string `json:"userPrincipalName,omitempty"`
}
type UserCollection struct {
	Users []User `json:"value"`
}
