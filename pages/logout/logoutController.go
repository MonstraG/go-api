package logout

import (
	"go-server/setup/myJwt"
	"go-server/setup/reqRes"
	"net/http"
	"time"
)

func GetHandler(w reqRes.MyWriter, _ *reqRes.MyRequest) {
	cookie := http.Cookie{
		Name:     myJwt.Cookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	w.RedirectToLogin()
}
