package app

import "net/http"

func Run() error {
	return http.ListenAndServe(BaseUrl, router)
}
