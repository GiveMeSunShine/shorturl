/**
 * @Author : ysh
 * @Description :
 * @File : middleware
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 下午14:31
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"shorturl/service"
	"time"
)

type dataResponse struct {
	LongUrl string `json:"long_url"`
	Code string `json:"code"`
	CreateAt time.Time `json:"create_at"`
	ShortUrl string `json:"short_url"`
}

type GetRequest struct {
	Code string
}

type GetResponse struct {
	Err error `json:"err"`
	Data interface{} `json:"data"`
}

type PostRequest struct {
	LongUrl string `json:"long_url" `
}

type PostResponse struct {
	Err error `json:"err"`
	Data dataResponse `json:"data"`
}

type Endpoints struct {
	GetEndpoint endpoint.Endpoint
	PostEndpoint endpoint.Endpoint
}


func MakeGetEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRequest)
		redirect, err := s.Get(ctx, req.Code)
		return GetResponse{Err: err,Data: redirect},err
	}
}

func MakePostEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostRequest)
		post, err := s.Post(ctx, req.LongUrl)
		resp := dataResponse{}

		if err == nil && post != nil {
			resp.Code = post.Code
			resp.LongUrl = req.LongUrl
			resp.CreateAt = post.CreatedAt
			resp.ShortUrl = post.ShortUrl
		}
		return PostResponse{Err: err,Data: resp},nil
	}
}

func (r GetResponse) Failed() error {
	return r.Err
}

type Failed interface {
	Failed() error
}


func New(s service.Service, mid map[string][] endpoint.Middleware) Endpoints {
	endpoints := Endpoints{
		GetEndpoint:  MakeGetEndpoint(s),
		PostEndpoint: MakePostEndpoint(s),
	}

	for _,m := range mid["Get"]{
		endpoints.GetEndpoint = m(endpoints.GetEndpoint)
	}

	for _,mid :=range mid["Post"]{
		endpoints.PostEndpoint = mid(endpoints.PostEndpoint)
	}
	return endpoints
}

func (e Endpoints) Get(ctx context.Context,code string) (rs interface{}, err error)  {
	request := GetRequest{Code: code}
	response, err := e.GetEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetResponse).Data,response.(GetResponse).Err
}

func (e Endpoints) Post(ctx context.Context, url string) (rs interface{}, err error) {
	request := PostRequest{
		LongUrl: url,
	}
	response, err := e.PostEndpoint(ctx, request)
	if err!=nil {
		return
	}
	return response.(PostResponse).Data,response.(PostResponse).Err
}
