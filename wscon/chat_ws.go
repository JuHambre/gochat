package wscon

import (
	"code.google.com/p/go.net/websocket"
	"crypto/cipher"
	"crypto/rc4"
	"crypto/sha256"
	"github.com/kobeld/gochatting/libs"
	"strings"
	"time"
)

const (
	TEXT_MTYPE   = "text_mtype"
	STATUS_MTYPE = "status_mtype"
	TIME_FORMAT  = "01-02 15:04:05"
)

var runningActiveRoom *ActiveRoom = &ActiveRoom{}

// Declaramos la clave para poder utilizarla tanto al enviar como recibir mensajes
var clave = "claveParaCifrar"

// Funcion que sirve para comprobar si hay error y salir del programa (panic) en tal caso
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func BuildConnection(ws *websocket.Conn) {
	email := ws.Request().URL.Query().Get("email")

	if email == "" {
		return
	}

	onlineUser := &OnlineUser{
		InRoom:     runningActiveRoom,
		Connection: ws,
		Send:       make(chan Message, 256),
		UserInfo: &User{
			Email:    email,
			Name:     strings.Split(email, "@")[0],
			Gravatar: libs.UrlSize(email, 20),
		},
	}
	runningActiveRoom.OnlineUsers[email] = onlineUser

	m := Message{
		MType: STATUS_MTYPE,
		UserStatus: UserStatus{
			Users: runningActiveRoom.GetOnlineUsers(),
		},
	}
	runningActiveRoom.Broadcast <- m

	go onlineUser.PushToClient()
	onlineUser.PullFromClient()

	onlineUser.killUserResource()
}

type ActiveRoom struct {
	OnlineUsers map[string]*OnlineUser
	Broadcast   chan Message
	CloseSign   chan bool
}

type OnlineUser struct {
	InRoom     *ActiveRoom
	Connection *websocket.Conn
	UserInfo   *User
	Send       chan Message
}

type User struct {
	Name     string
	Email    string
	Gravatar string
}

type Message struct {
	MType       string
	TextMessage TextMessage
	UserStatus  UserStatus
}

type TextMessage struct {
	Content  string
	UserInfo *User
	Time     string
}

type UserStatus struct {
	Users []*User
}

func InitChatRoom() {
	runningActiveRoom = &ActiveRoom{
		OnlineUsers: make(map[string]*OnlineUser),
		Broadcast:   make(chan Message),
		CloseSign:   make(chan bool),
	}
	go runningActiveRoom.run()
}

// Core function of room
func (this *ActiveRoom) run() {
	for {
		select {
		case b := <-this.Broadcast:
			for _, online := range this.OnlineUsers {
				online.Send <- b
			}
		case c := <-this.CloseSign:
			if c == true {
				close(this.Broadcast)
				close(this.CloseSign)
				return
			}
		}
	}
}

func (this *OnlineUser) PullFromClient() {
	for {
		var content string
		err := websocket.Message.Receive(this.Connection, &content)
		// If user closes or refreshes the browser, a err will occur
		if err != nil {
			return
		}

		// Podemos modificar el texto que reciben todos los usuarios
		// content = "prueba"

		// Obtenemos hash de 256 con la clave
		h := sha256.New()
		h.Reset()
		h.Write([]byte(clave))
		key := h.Sum(nil)

		// Utilizamos la interfaz cipher para descifrar y cifrar
		var S cipher.Stream

		// RC4 proporciona un cipher.Stream directamente al ser cifrado en flujo
		a, err := rc4.NewCipher(key)
		check(err)
		S = a

		// Desciframos y parseamos de []byte a string
		cadCifrada := []byte(content)
		plaintext := make([]byte, len(cadCifrada))
		S.XORKeyStream(plaintext, cadCifrada)
		content = string(plaintext[:])

		// Mostramos el mensaje a todos

		m := Message{
			MType: TEXT_MTYPE,
			TextMessage: TextMessage{
				UserInfo: this.UserInfo,
				Time:     humanCreatedAt(),
				Content:  content,
			},
		}
		this.InRoom.Broadcast <- m
	}
}

func (this *OnlineUser) PushToClient() {
	for b := range this.Send {

		// Podemos modificar el texto que se envia
		// b.TextMessage.Content = "prueba"

		// Obtenemos hash de 256 con la clave
		h := sha256.New()
		h.Reset()
		h.Write([]byte(clave))
		key := h.Sum(nil)

		// Utilizamos la interfaz cipher para descifrar y cifrar
		var S cipher.Stream

		// RC4 proporciona un cipher.Stream directamente al ser cifrado en flujo
		c, error := rc4.NewCipher(key)
		check(error)
		S = c

		// Ciframos
		plaintext := []byte(b.TextMessage.Content)
		cadCifrada := make([]byte, len(plaintext))
		S.XORKeyStream(cadCifrada, plaintext)

		// Enviamos parseando de []byte a string
		b.TextMessage.Content = string(cadCifrada[:])

		// Enviamos la cadena cifrada

		err := websocket.JSON.Send(this.Connection, b)
		if err != nil {
			break
		}
	}
}

func (this *OnlineUser) killUserResource() {
	this.Connection.Close()
	delete(this.InRoom.OnlineUsers, this.UserInfo.Email)
	close(this.Send)

	m := Message{
		MType: STATUS_MTYPE,
		UserStatus: UserStatus{
			Users: runningActiveRoom.GetOnlineUsers(),
		},
	}
	runningActiveRoom.Broadcast <- m
}

func (this *ActiveRoom) GetOnlineUsers() (users []*User) {
	for _, online := range this.OnlineUsers {
		users = append(users, online.UserInfo)
	}
	return
}

func humanCreatedAt() string {
	return time.Now().Format(TIME_FORMAT)
}
