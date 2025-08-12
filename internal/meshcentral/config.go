package meshcentral

import (
	"time"

	"github.com/gorilla/websocket"
)

type Device struct {
	Id   string
	Name string
	OS   string
	IP   string
	Icon int
	Conn int
	Pwr  int
}

type NewTokenResponse struct {
	Name      string
	TokenUser string
	TokenPass string
	Created   int64
	Expire    int64
}

type Settings struct {
	ServerURL        string
	Username         string
	Password         string
	AuthCookie       string
	ServerID         string
	LoginKey         string
	LocalPort        int
	RemotePort       int
	RemoteTarget     string
	RemoteNodeID     string
	WebSocket        *websocket.Conn
	WebChannel       chan struct{}
	ACookie          string
	RCookie          string
	RenewCookieTimer *time.Timer
	Devices          []Device
	DeviceQueryState int
	TokenQueryState  int
	TokenResponse    NewTokenResponse
	debug            bool
	MfaToken         string
	//EmailToken     bool
	//SMSToken       bool
}

var settings Settings

func ApplySettings(remoteNodeId string, remotePort int, localPort int, remoteTarget string, debug bool) {
	settings.RemoteNodeID = remoteNodeId
	settings.RemotePort = remotePort
	settings.LocalPort = localPort
	settings.RemoteTarget = remoteTarget
	settings.debug = debug
}

func SetLoginInfo(username, password, server string) {
	settings.Username = username
	settings.Password = password
	settings.ServerURL = "wss://" + server + "/meshrelay.ashx"
}

func SetMfaToken(token string) {
	settings.MfaToken = token
}
