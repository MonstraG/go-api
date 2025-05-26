package logout

import (
	"go-api/infrastructure/reqRes"
)

func GetHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	w.ExpireCookie()
	w.RedirectToLogin(r)
}
