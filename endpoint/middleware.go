/**
 * @Author : ysh
 * @Description :
 * @File : middleware
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 上午11:31
 */

package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"golang.org/x/time/rate"
	"time"
)

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint{
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				_ = logger.Log("transport_err",err,"took",time.Since(begin))
			}(time.Now())
			return next(ctx,request)
		}
	}
}

var LimitExeedErr = errors.New("Rate limit exceed")

func NewTokenBucketLimitter(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow(){
				return nil, LimitExeedErr
			}
			return next(ctx,request)
		}
	}
}
