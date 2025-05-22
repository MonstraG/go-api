package logout

import (
	"go-server/setup/reqRes"
)

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	w.ExpireCookie()
	w.RedirectToLogin(r)
}
