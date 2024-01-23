package models

type ResponseBody struct {
	Username     string `json:"username"`
	AccesToken   string `json:"access_token"`
	RefreshToken string `json:"refesh_token"`
}
