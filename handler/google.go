package handler

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"simka/config"
	"simka/utilities"
)

type googleHandler struct{}

var Token oauth2.Token

func NewGoogleHandler() *googleHandler {
	return &googleHandler{}
}

func (h *googleHandler) GetGoogleLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login-google.html", gin.H{"content": "ok"})
}

func (h *googleHandler) PostGoogleLogin(c *gin.Context) {
	oauthState := utilities.OauthState

	u := config.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (h *googleHandler) GetGoogleCallBack(c *gin.Context) {
	oauthState := utilities.OauthState
	state := c.Query("state")
	code := c.Query("code")

	if state != oauthState {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	token, err := config.AppConfig.GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		h.Error(err.Error(), c)
		return
	}

	client := http.Client{}
	resp, err := client.Get(config.OauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		h.Error(err.Error(), c)
		return
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.Error(err.Error(), c)
		return
	}
	_ = contents

	Token = *token
	c.Redirect(http.StatusTemporaryRedirect, "/simaka")
}

func (s *googleHandler) Error(errMsg string, c *gin.Context) {
	c.HTML(http.StatusBadRequest, "error.html", gin.H{"content": errMsg})
}

func (s *googleHandler) Close(c *gin.Context) {
	Token = oauth2.Token{}
}
