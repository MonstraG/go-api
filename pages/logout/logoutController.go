package logout

import (
	"go-server/setup/reqRes"
)

func GetHandler(w reqRes.MyWriter, _ *reqRes.MyRequest) {
	w.ExpireCookie()
	w.RedirectToLogin()
}
