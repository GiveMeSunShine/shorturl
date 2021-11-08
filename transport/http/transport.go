/**
 * @Author : ysh
 * @Description : http 请求参数编码，解码
 * @File : transport
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 上午11:44
 */

package http

import (
	"context"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"shorturl/endpoint"
	"shorturl/service"
)

func NewHttpHandler(endpoints endpoint.Endpoints, options map[string][] kithttp.ServerOption) http.Handler {
	router := mux.NewRouter()
	router.Handle("/{code}",kithttp.NewServer(
			endpoints.GetEndpoint,
			decodeGetRequest,
			encodeGetResponse,
			options["Get"]...
		)).Methods(http.MethodGet)

	router.Handle("/",kithttp.NewServer(
			endpoints.PostEndpoint,
			decodePostRequest,
			encodePostResponse,
			options["Post"]...
		)).Methods(http.MethodPost)
	return router
}

func decodePostRequest(_ context.Context,r *http.Request) (interface{},error) {
	var req endpoint.PostRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b,&req); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err:= validate.Struct(req); err != nil {
		return nil, errors.Wrap(err, service.DataInvalidErr.Error())
	}
	return req,nil
}

func encodePostResponse(ctx context.Context, w http.ResponseWriter,response interface{}) (err error) {
	if f,ok := response.(endpoint.Failed); ok && f.Failed() != nil{
		ErrorEncoder(ctx,f.Failed(),w)
		return nil
	}
	err = json.NewEncoder(w).Encode(response)
	return
}

var (
	ErrCodeNotFound = errors.New("code is nil")
)

func decodeGetRequest(_ context.Context, request *http.Request) (interface{}, error) {
	vars := mux.Vars(request)
	code , ok := vars["code"]
	if !ok {
		return nil,ErrCodeNotFound
	}
	req := endpoint.GetRequest{
		Code: code,
	}
	return req,nil
}

func encodeGetResponse(ctx context.Context, w http.ResponseWriter,response interface{}) ( err error) {
	if f,ok := response.(endpoint.Failed); ok && f.Failed() != nil {
		ErrorRedirect(ctx,f.Failed(),w)
		return nil
	}
	resp := response.(endpoint.GetResponse)
	redirect := resp.Data.(*service.Redirect)
	http.Redirect(w,&http.Request{},redirect.LongUrl,http.StatusFound)
	return
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func ErrorRedirect(_ context.Context, err error, w http.ResponseWriter) {
	http.Redirect(w, &http.Request{}, os.Getenv("SHORT_URI"), http.StatusFound)
}

func err2code(err error) int {
	return http.StatusOK
}

type errorWrapper struct {
	Error string `json:"error"`
}



