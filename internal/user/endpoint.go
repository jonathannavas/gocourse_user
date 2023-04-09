package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jonathannavas/gocourse_meta/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	UpdateRequest struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimitPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var userBody CreateRequest

		err := json.NewDecoder(r.Body).Decode(&userBody)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Invalid request format"})
			return
		}

		if userBody.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Firstname is required"})
			return
		}

		if userBody.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Lastname is required"})
			return
		}

		user, err := s.Create(userBody.FirstName, userBody.LastName, userBody.Email, userBody.Phone)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 201, Data: user})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		user, err := s.Get(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		queryParams := r.URL.Query()
		filters := Filters{
			FirstName: queryParams.Get("first_name"),
			LastName:  queryParams.Get("last_name"),
		}

		limit, _ := strconv.Atoi(queryParams.Get("limit"))
		page, _ := strconv.Atoi(queryParams.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Error: err.Error()})
			return
		}

		meta, err := meta.New(page, limit, count, config.LimitPageDef)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Error: err.Error()})
			return
		}

		users, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: users, Meta: meta})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("update users service")

		var body UpdateRequest

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Invalid request format"})
			return
		}

		if body.FirstName != nil && *body.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Firstname is required"})
			return
		}

		if body.LastName != nil && *body.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Lastname is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Update(id, body.FirstName, body.LastName, body.Email, body.Phone); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Error: "User doesn't exists"})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		err := s.Delete(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}
