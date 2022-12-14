package simka

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"

	"simka/helper"
	"simka/utilities"
)

type Service interface {
	ParseWebLogin(input UserLoginInput) (UserLoginInput, error)
	FormatEvent(events helper.ListMataKuliah) ([]EventFormatter, error)
	CreateEvent(token oauth2.Token, event []EventFormatter) error
	Login(input UserLoginInput) (UserLoginInput, error)
	ParseMataKuliah() (helper.ListMataKuliah, error)
	SetEndTime(date int)
}

type service struct {
	Jar          *utilities.Jar
	Client       http.Client
	SuccessLogin bool
	EndMonthDate int
}

var (
	urlMaktul = "https://simaka.asia.ac.id/mahasiswa/dashboard.php?p=LJgsnzSxq2Sf&idmenu=AD==&idsubmenu=ZGN="
	urlPost   = "https://simaka.asia.ac.id/otentikasi.php"
	urlIndex  = "https://simaka.asia.ac.id/index.php"
)

func NewService(jar *utilities.Jar) *service {
	return &service{Jar: jar}
}

func (s *service) SetEndTime(date int) {
	s.EndMonthDate = date
}

func (s *service) Login(input UserLoginInput) (UserLoginInput, error) {
	errUser := "Username dan Password tidak Valid"

	userData, err := s.ParseWebLogin(input)
	if err != nil {
		return input, err
	}

	cookieJar := s.Jar
	s.Client = http.Client{Jar: cookieJar}

	formData := url.Values{
		"username": {userData.NIM},
		"password": {userData.Password},
		"jawaban":  {userData.UserAnswer},
		"benar":    {userData.SystemAnswer},
	}

	resp, err := s.Client.PostForm(urlPost, formData)
	if err != nil {
		return input, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return input, err
	}

	if strings.Contains(string(bodyData), errUser) {
		return input, errors.New(errUser)
	}

	s.SuccessLogin = true
	return input, nil
}

func (s *service) ParseWebLogin(input UserLoginInput) (UserLoginInput, error) {
	answerQuestion, err := helper.ParseAnswerQuestion(urlIndex)
	if err != nil {
		return input, err
	}

	input.UserAnswer = answerQuestion
	input.SystemAnswer = answerQuestion

	return input, nil
}

func (s *service) ParseMataKuliah() (helper.ListMataKuliah, error) {
	var listMatkul helper.ListMataKuliah
	errMsg := "Need login first"

	if status := s.SuccessLogin; !status {
		return listMatkul, errors.New(errMsg)
	}

	resp, err := s.Client.Get(urlMaktul)
	if err != nil {
		return listMatkul, errors.New(errMsg)
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return listMatkul, err
	}

	sourceHtmlMatkul := helper.ParseDataTableMatkul(string(bodyData))
	listMatkul = helper.StoreDataMatkul(sourceHtmlMatkul)

	return listMatkul, nil
}

func (s *service) FormatEvent(events helper.ListMataKuliah) ([]EventFormatter, error) {
	eventFormatter, err := CalendarEventFormatter(events)
	if err != nil {
		return eventFormatter, err
	}

	return eventFormatter, nil
}

func (s *service) CreateExpiredMonthDate(timeNow time.Time) []string {
	tnNext := timeNow.AddDate(0, s.EndMonthDate-1, 0)
	monthNow := int(timeNow.Month())
	tnMonthNext := int(tnNext.Month())

	index := 1
	overFlow := false

	listMonth := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	listTargetMonth := make([]string, 0)
	hashMap := make(map[string]bool, 0)

	appendMonth := func(index int) {
		nMonth := listMonth[index]
		if nMonth == 13 {
			nMonth = 1
		}

		for _, val := range listTargetMonth {
			hashMap[val] = true
		}

		if isExits := hashMap[strconv.Itoa(nMonth)]; !isExits {
			listTargetMonth = append(listTargetMonth, strconv.Itoa(nMonth))
		}
	}

	for true {
		if index == 13 {
			overFlow = true
			index = 1
		}

		if index >= monthNow && !overFlow {
			appendMonth(index)
			if index == tnMonthNext {
				break
			}
		}

		if overFlow {
			appendMonth(index)
			if index == tnMonthNext {
				break
			}
		}

		index += 1
	}

	return listTargetMonth
}

func (s *service) CreateEvent(token oauth2.Token, events []EventFormatter) error {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&token))

	calendarService, err := calendar.New(client)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	timeNow := time.Now()
	listTargetMonth := s.CreateExpiredMonthDate(timeNow)

	for _, event := range events {
		tn := timeNow.AddDate(0, 0, event.HtoDay)

		year := tn.Year()
		month := tn.Month()
		day := tn.Day()
		sec := tn.Second()
		msec := tn.Nanosecond()

		startTime := time.Date(year, month, day, event.StartHour, event.StartMin, sec, msec, loc).Format(time.RFC3339)
		endTime := time.Date(year, month, day, event.EndHour, event.EndMin, sec, msec, loc).Format(time.RFC3339)
		recurrence := fmt.Sprintf("%s%s", event.Recurrence, strings.Join(listTargetMonth, ","))
		newEvent := calendar.Event{
			Summary:     event.Summary,
			Description: event.Description,
			Start: &calendar.EventDateTime{
				DateTime: startTime,
				TimeZone: loc.String(),
			},
			End: &calendar.EventDateTime{
				DateTime: endTime,
				TimeZone: loc.String(),
			},
			Recurrence: []string{recurrence},
		}

		_, err := calendarService.Events.Insert("primary", &newEvent).Do()
		if err != nil {
			return err
		}
	}

	return nil
}
