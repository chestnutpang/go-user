package transport

/*
transport 包中，我们需将构建好的 Endpoint 通过 HTTP 或者 RPC 方式暴露
*/

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/longjoy/micro-go-course/endpoint"
	"net/http"
	"os"
)
var (
	ErrorBadRequest = errors.New("invalid request parameter")
)
// 使用 mux 构建 http 处理
func MakeHttpHandler(ctx context.Context, endpoints *endpoint.UserEndpoints) http.Handler{
	r := mux.NewRouter()

	kitLog := log.NewLogfmtLogger(os.Stderr)
	kitLog = log.With(kitLog, "ts", log.DefaultTimestampUTC)
	kitLog = log.With(kitLog, "caller", log.DefaultCaller)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	// register 路由处理
	r.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeJSONResponse,
		options...,
	))

	return r
}


// 读取 http 请求中发送的数据，并处理封装成 LoginRequest 请求体
func decodeLoginRequest(_ context.Context, r *http.Request)(interface{}, error){
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == ""{
		return nil, ErrorBadRequest
	}
	return &endpoint.LoginRequest{
		Email: email,
		Password: password,
	}, nil
}

func decodeRegisterRequest(_ context.Context, r *http.Request)(interface{}, error){
	fmt.Println(r.Form)
	fmt.Println(r.PostForm)

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	//username := r.PostFormValue("username")
	//email := r.PostFormValue("email")
	//password := r.PostFormValue("password")
	if email == "" || password == "" || username == ""{
		return nil, ErrorBadRequest
	}
	
	return &endpoint.RegisterRequest{
		Username: username,
		Password: password,
		Email: email,
	}, nil
}


// 对json响应进行编码
func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error{
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}


func encodeError(_ context.Context, err error, w http.ResponseWriter){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err{
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}