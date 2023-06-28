package dto

type SigninValues struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
type SigninResult struct {
	Succeeded    bool   `json:"succeeded"`
	ErrorMessage string `json:"errorMessage"`
	NextURL      string `json:"nextUrl"`
}
