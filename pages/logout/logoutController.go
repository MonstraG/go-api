package logout

import (
	"go-api/infrastructure/reqRes"
)

func PerformLogout(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	w.ExpireCookie()
	w.RedirectToLogin(r)
}
