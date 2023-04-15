package user

import (
	"context"
	"errors"

	"github.com/jonathannavas/go_lib_response/response"
	"github.com/jonathannavas/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)
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

	GetReq struct {
		ID string
	}

	DeleteReq struct {
		ID string
	}

	GetAllReq struct {
		FirstName string
		LastName  string
		Limit     int
		Page      int
	}

	UpdateRequest struct {
		ID        string
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		userBody := request.(CreateRequest)

		if userBody.FirstName == "" {
			return nil, response.BadRequest(errFirstnameRequired.Error())
		}

		if userBody.LastName == "" {
			return nil, response.BadRequest(errLastnameRequired.Error())
		}

		user, err := s.Create(ctx, userBody.FirstName, userBody.LastName, userBody.Email, userBody.Phone)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil, 200), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		user, err := s.Get(ctx, req.ID)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", user, nil, 201), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimitPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", users, meta, 201), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UpdateRequest)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(errFirstnameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(errLastnameRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil, 201), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		err := s.Delete(ctx, req.ID)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil, 200), nil
	}
}
