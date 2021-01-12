package main

import (
	"net/http"
	"context"
)

func (a *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch r.URL.Path {
	case "/user/profile":
		a.Create(ctx, CreateParams{Login: "hello",Age: 10,Name:"sda",Status:"moderator"})
	}
}