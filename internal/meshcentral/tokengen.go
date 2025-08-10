package meshcentral

import (
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func NewLoginToken(name string, expire int) NewTokenResponse {
	settings.TokenQueryState = 1
	if settings.debug {
		fmt.Println("Sending createLoginToken command")
	}

	// send createLoginToken message
	msg := fmt.Sprintf(`{"action":"createLoginToken","name":"%s","expire":%d}`, name, expire)
	settings.WebSocket.WriteMessage(websocket.TextMessage, []byte(msg))

	// wait for TokenQueryState to change
	for settings.TokenQueryState == 1 {
		time.Sleep(100 * time.Millisecond)
	}

	return settings.TokenResponse
}

func handleNewTokenCommand(command map[string]interface{}) {
	if settings.debug {
		fmt.Println("Received createLoginToken command")
	}
	var response NewTokenResponse

	response.Name = command["name"].(string)
	response.TokenUser = command["tokenUser"].(string)
	response.TokenPass = command["tokenPass"].(string)
	response.Created = int64(command["created"].(float64))
	response.Expire = int64(command["expire"].(float64))

	settings.TokenResponse = response
	settings.TokenQueryState = 0
}

func GenerateLoginToken(server string, username string, password string) (tokenUser string, tokenPass string, err error) {
	ApplySettings(
		"",
		0,
		0,
		"",
		false,
	)

	SetLoginInfo(username, password, server)
	StartSocket()

	// create login token name as mcc-<hostname>-<yyyymmdd>
	hostname, _ := os.Hostname()
	tokenName := "mcc-" + hostname + "-" + time.Now().Format("20060102")

	token := NewLoginToken(tokenName, 0)
	StopSocket()

	return token.TokenUser, token.TokenPass, nil
}
