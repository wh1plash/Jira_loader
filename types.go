package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type GetTokenRequest struct {
	UserName string `json:"username"`
	Pass     string `json:"password"`
	Project  string `json:"project"`
}

type LoaderToken struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
	Exp   int64  `json:"exp"`
}

type TokenClaims struct {
	Exp int64 `json:"exp"`
	jwt.RegisteredClaims
}

type LoadOrderResponse struct {
	TotalRow int `json:"totalRow"`
	SavedRow int `json:"savedRow"`
	ErrorRow int `json:"errorRow"`
}

type JiraComment struct {
	Body string `json:"body"`
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
	token   LoaderToken
}
