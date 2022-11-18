package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simka/simka"
)

type SimkaHandler struct {
	simakaService simka.Service
}

func NewSimakaHandler(simakaSvc simka.Service) *SimkaHandler {
	return &SimkaHandler{simakaSvc}
}

func (s *SimkaHandler) GetLoginSimaka(c *gin.Context) {
	c.HTML(http.StatusOK, "login-simaka.html", gin.H{})
}

func (s *SimkaHandler) PostLoginSimaka(c *gin.Context) {
	var input simka.UserLoginInput

	nim := c.PostForm("nim")
	password := c.PostForm("password")

	input.NIM = nim
	input.Password = password

	_, err := s.simakaService.Login(input)
	if err != nil {
		s.Error("Invalid username or password", c)
		return
	}
	s.GetListMataKuliah(c)
}

func (s *SimkaHandler) GetListMataKuliah(c *gin.Context) {
	token := Token

	data, err := s.simakaService.ParseMataKuliah()
	if err != nil {
		s.Error(err.Error(), c)
		return
	}
	dataFormatter, err := s.simakaService.FormatEvent(data)
	if err != nil {
		s.Error(err.Error(), c)
		return
	}

	err = s.simakaService.CreateEvent(token, dataFormatter)
	if err != nil {
		s.Error(err.Error(), c)
		return
	}

	// redirect to success page
	s.SuccessLogin(c)
}

func (s *SimkaHandler) SuccessLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "success.html", gin.H{})
}

func (s *SimkaHandler) Error(errorMsg string, c *gin.Context) {
	c.HTML(http.StatusBadRequest, "error.html", gin.H{"content": errorMsg})
}
