package utilities

import (
	"crypto/rand"
	"encoding/base64"
)

var OauthState string

func GenerateOauthCookie() {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)
	OauthState = state
}
