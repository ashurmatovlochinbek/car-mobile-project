package user

import "net/http"

type Handlers interface {
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	Refresh() http.HandlerFunc
	VerifyOTP() http.HandlerFunc
	GetSecuredResource() http.HandlerFunc
}
