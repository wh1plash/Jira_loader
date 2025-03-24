package main

import (
	"net/http"
	"time"
)

type GetTokenRequest struct {
	UserName string `json:"username"`
	Pass     string `json:"password"`
	Project  string `json:"project"`
}

type GetTokenResponse struct {
	Token string `json:"access_token"`
	TTL   int    `json:"expires_in"`
}

type LoadOrderRequest struct {
}

type LoadOrderResponse struct {
	TotalRow int `json:"totalRow"`
	SavedRow int `json:"savedRow"`
	ErrorRow int `json:"errorRow"`
}

type Status struct {
	Name string `json:"name"`
}

type Fields struct {
	Summary string `json:"summary,omitempty"`
	Status  `json:"status"`
}

type Queue struct {
	Values []struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
		Self   string `json:"self"`
		Fields `json:"fields"`
	} `json:"values"`
}

type Task struct {
	ID     string `json:"id"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		Description string `json:"description"`
		Attachment  []struct {
			FileName string `json:"filename"`
			Content  string `json:"content"`
		} `json:"attachment"`
	} `json:"fields"`
}

type HTTPClient struct {
	client  *http.Client
	baseUrl string
}

type Order struct {
	Type      string
	Count     int
	Region    int
	ProjectID int
	SerialNum int
	DateTo    time.Time
	Serie     int
}
