package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func contains(arr []string, str string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}
	return false
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	d, _ := json.Marshal(data)
	w.Write(d)
}

func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// endpoints
	switch r.URL.Path {
	case "/user/create":
		h.wrapperCreate(w, r)
	}
}

func parseCreateParams(r *http.Request) (CreateParams, *ApiError) {
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
		return params, &ApiError{Err: fmt.Errorf("invalid fields"), HTTPStatus: http.StatusBadRequest}
	}

	return params, nil
}

func (h *MyApi) wrapperCreate(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var params CreateParams
	if p, err := parseCreateParams(r); err == nil {
		params = p
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.Create(ctx, params)
	if err == nil {
		// send success response
		w.WriteHeader(http.StatusOK)
		respondJSON(w, res)
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
