package player

import (
	"go-api/infrastructure/reqRes"
	"html/template"
)

var playerTemplate = template.Must(template.ParseFiles("pages/player/playerPartial.gohtml"))

func GetPlayer(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	w.RenderTemplate(playerTemplate, r)
}
