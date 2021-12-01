/**
 * @Author : ysh
 * @Description :
 * @File : middleware
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 上午10:41
 */

package service

import (
	"context"
	"github.com/go-kit/log"
)

type Middleware func(svc Service) Service

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (l loggingMiddleware) Get(ctx context.Context, code string) (redirect *Redirect, err error) {
	defer func() {
		_ = l.logger.Log("method", "Get", "s", code, "err", err)
	}()
	return l.next.Get(ctx, code)
}

func (l loggingMiddleware) Post(ctx context.Context, info PostInfo) (redirect *Redirect, err error) {
	defer func() {
		_ = l.logger.Log("method", "Post", "url", info, "err", err)
	}()
	return l.next.Post(ctx, info)
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(svc Service) Service {
		return &loggingMiddleware{logger, svc}
	}
}
