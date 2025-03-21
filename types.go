package main

import (
	"net/http"
	"time"
)

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
		Attachment []struct {
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
