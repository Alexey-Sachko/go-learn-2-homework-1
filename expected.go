package main

import (
	"net/http"
	"context"
	"strconv"
	"encoding/json"
	"fmt"
)

func contains(arr []string, str string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}
	return false
}

func parseAndValidate(r *http.Request) (CreateParams, *ApiError) {
	params := CreateParams{}

	errorFields := []string{}

	if login := r.FormValue("Login"); login != "" {
		params.Login = login
	} else {
		errorFields = append(errorFields, "Login")
	}


	name := r.FormValue("full_name")
	params.Name = name
	
	statusEnum := []string{"user", "admin", "moderator"}
	if status := r.FormValue("Status"); contains(statusEnum, status) {
		params.Status = status
	} else {
		errorFields = append(errorFields, "Status")
	}
	

	if a, err := strconv.Atoi(r.FormValue("Age")); err == nil {
		params.Age = a
	} else {
		errorFields = append(errorFields, "Age")
	}

	if len(errorFields) > 0 {
		return params, &ApiError{Err: fmt.Errorf("invalid fields"),HTTPStatus: http.StatusBadRequest}
	}

	return params, nil
}

func (a *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch r.URL.Path {
	case "/user/create":
		params := CreateParams{}

		res, err := a.Create(ctx, params)
		if err == nil {
			// send success response
			w.WriteHeader(http.StatusOK)
			d, _ := json.Marshal(res) // TODO
			w.Write(d)
			return
		}

		switch err.(type) {
		case ApiError:
			e := err.(ApiError)
			w.WriteHeader(e.HTTPStatus)
			w.Write([]byte(e.Error()))
		
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}