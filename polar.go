package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type trainingHistoryRequest struct {
	FromDate string `json:"fromDate"`
	SportIds []int  `json:"sportIds"`
	ToDate   string `json:"toDate"`
	UserID   int    `json:"userId"`
}

type userData struct {
	ID int `json:"id"`
}

type currentUserResponse struct {
	User userData `json:"user"`
}

type Session struct {
	Jar    *cookiejar.Jar
	Client http.Client
	User   userData
}

const (
	Meter     = 1
	Kilometer = 1000 * Meter
)

type Distance float32

func (d Distance) String() string {
	km := int(d) / Kilometer
	m := int(d) - km*Kilometer
	return fmt.Sprintf("%d.%dkm", km, m)
}

type Training struct {
	ID        int           `json:"id"`
	Duration  time.Duration `json:"duration"`
	Distance  Distance      `json:"distance"`
	StartDate string        `json:"startDate"`
}

func NewSession(username string, password string) Session {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	resp, _ := client.PostForm("https://flow.polar.com/login",
		url.Values{"email": {username}, "password": {password}})
	resp.Body.Close()

	return Session{Jar: jar, Client: client}
}

func (s Session) UpdateUserData() {
	s.User = s.GetUserData()
}

func (s Session) GetUserData() userData {
	resp, _ := s.Client.Get("https://flow.polar.com/api/account/users/current/user")
	var data currentUserResponse
	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	return data.User
}

func (s Session) GetTrainings(start string, end string, sportIDs []int) []Training {
	request := &trainingHistoryRequest{
		FromDate: start,
		ToDate:   end,
		SportIds: sportIDs,
		UserID:   s.User.ID}
	v, _ := json.Marshal(request)

	resp, _ := s.Client.Post("https://flow.polar.com/api/training/history", "application/json", bytes.NewBuffer(v))
	var trainings []Training
	err := json.NewDecoder(resp.Body).Decode(&trainings)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(trainings); i++ {
		trainings[i].Duration *= time.Millisecond
	}

	return trainings
}

func (s Session) GetTrainingGpx(trainingId int) []byte {
	resp, _ := s.Client.Get(fmt.Sprintf("https://flow.polar.com/api/export/training/gpx/%d", trainingId))
	data, _ := ioutil.ReadAll(resp.Body)
	return data
}
