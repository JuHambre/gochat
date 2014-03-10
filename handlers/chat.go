package handlers

import (
	. "github.com/paulbellamy/mango"
	"github.com/sunfmin/mangotemplate"
	"net/http"
)

type RenderData struct {
	Email         string
	WebSocketHost string
}

func Home(env Env) (status Status, headers Headers, body Body) {
	mangotemplate.ForRender(env, "chats/home", nil)
	headers = Headers{}
	return
}

func Join(env Env) (status Status, headers Headers, body Body) {
	email := env.Request().FormValue("email")
	//Anyadimos contrasenya
	password := env.Request().FormValue("password")
	//Comprobamos ahora tambien la contrasenya
	if email == "" || password == "" {
		return Redirect(http.StatusFound, "/")
	}

	r := env.Request()
	mangotemplate.ForRender(env, "chats/room", &RenderData{Email: email, WebSocketHost: r.Host})
	headers = Headers{}
	return
}
