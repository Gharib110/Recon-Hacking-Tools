package rpc

import (
	"bytes"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
	"net/http"
)

type Metasploit struct {
	host  string
	user  string
	pass  string
	token string
}

type loginRequest struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Username string
	Password string
}

type LoginResponse struct {
	Result       string `msgpack:"result"`
	Token        string `msgpack:"token"`
	Error        bool   `msgpack:"error"`
	ErrorClass   string `msgpack:"error_class"`
	ErrorMessage string `msgpack:"error_message"`
}

type logoutRequest struct {
	_msgpack    struct{} `msgpack:",asArray"`
	Method      string
	Token       string
	LogoutToken string
}

type LogoutResponse struct {
	Result string `msgpack:"result"`
}

type sessionListRequest struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Token    string
}

type SessionListResponse struct {
	ID          uint32 `msgpack:",omitempty"`
	Type        string `msgpack:"type"`
	TunnelLocal string `msgpack:"tunnel_local"`
	TunnelPeer  string `msgpack:"tunnel_peer"`
	ViaPayload  bool   `msgpack:"via_payload"`
	Description string `msgpack:"description"`
	Info        string `msgpack:"info"`
	Workspace   string `msgpack:"workspace"`
	SessionHost string `msgpack:"session_host"`
	SessionPort string `msgpack:"session_port"`
	Username    string `msgpack:"username"`
	UUID        string `msgpack:"uuid"`
	ExploitUUID string `msgpack:"exploit_uuid"`
}

// New creates a new metasploit object, login then return it
func New(host string, user string, pass string, token string) (*Metasploit, error) {
	msf := &Metasploit{
		host:  host,
		user:  user,
		pass:  pass,
		token: token,
	}

	err := msf.Login()
	if err != nil {
		return nil, err
	}

	return msf, nil
}

// Login just login
func (msf *Metasploit) Login() error {
	ctx := &loginRequest{
		Method:   "auth.login",
		Username: msf.user,
		Password: msf.pass,
	}
	var res LoginResponse
	err := msf.Send(ctx, &res)
	if err != nil {
		return err
	}

	msf.token = res.Token
	return nil
}

func (msf *Metasploit) Send(ctx interface{}, l interface{}) error {
	buf := new(bytes.Buffer)
	msgpack.NewEncoder(buf).Encode(ctx)
	dest := fmt.Sprintf("http://%s/api", msf.host)
	r, err := http.Post(dest, "binary/message-pack", buf)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = msgpack.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		return err
	}

	return nil
}

func (msf *Metasploit) Logout() error {
	ctx := &logoutRequest{
		Method:      "auth,logout",
		Token:       msf.token,
		LogoutToken: msf.token,
	}

	var res LogoutResponse
	err := msf.Send(ctx, &res)
	if err != nil {
		return err
	}

	msf.token = ""
	return nil
}

func (msf *Metasploit) SessionList() (map[uint32]SessionListResponse, error) {
	req := &sessionListRequest{
		Method: "session.list",
		Token:  msf.token,
	}
	res := make(map[uint32]SessionListResponse)
	err := msf.Send(req, &res)
	if err != nil {
		return nil, err
	}

	for id, session := range res {
		session.ID = id
		res[id] = session
	}

	return res, nil
}
