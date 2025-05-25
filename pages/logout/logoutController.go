package logout

import (
	"go-api/setup/reqRes"
)

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	w.ExpireCookie()
	w.RedirectToLogin(r)
}
