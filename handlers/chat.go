package handlers

import (
	"bufio"
	. "github.com/paulbellamy/mango"
	"github.com/sunfmin/mangotemplate"
	"net/http"
	"os"
	"strings"
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

// Funcion que sirve para comprobar si hay error y salir del programa (panic) en tal caso
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Join(env Env) (status Status, headers Headers, body Body) {
	email := env.Request().FormValue("email")
	// Anyadimos contrasenya
	password := env.Request().FormValue("password")
	// Comprobamos ahora tambien la contrasenya
	if email == "" || password == "" {
		return Redirect(http.StatusFound, "/")
	}

	// Leemos desde fichero de texto y comprobamos si existe el usuario
	file, err := os.Open("usuarios.txt")
	check(err)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " ")
		if s[0] == email {
			if s[1] == password {
				r := env.Request()
				mangotemplate.ForRender(env, "chats/room", &RenderData{Email: email, WebSocketHost: r.Host})
				headers = Headers{}
				return
			}
		}
	}

	return Redirect(http.StatusFound, "/")

}
