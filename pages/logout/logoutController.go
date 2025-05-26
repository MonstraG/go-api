package logout

import (
	"go-api/infrastructure/reqRes"
)

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	w.ExpireCookie()
	w.RedirectToLogin(r)
}
