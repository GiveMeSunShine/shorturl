/**
 * @Author : ysh
 * @Description :
 * @File : service
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 上午10:33
 */

package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/log"
	"shorturl/shortid"
	"strings"
	"time"
)

var (
	DataNotFoundErr = errors.New("Redirect Data Not Found")
	DataInvalidErr = errors.New("Redirect Data Invalid")
)

type Service interface {
	Get(ctx context.Context, code string) (redirect *Redirect, err error)
	Post(ctx context.Context, demain string) (redirect *Redirect, err error)
}

type service struct {
	repository Repository
	logger log.Logger
	shortUrl string
	maxLen int
}

func (s service) Get(ctx context.Context, code string) (redirect *Redirect, err error) {
	return s.repository.Find(code)
}

func (s service) Post(ctx context.Context, url string) (redirect *Redirect, err error) {
	now := time.Now()
	now = now.In(time.Local)
	var code string
	code = shortid.ShortIdFromCode()
	if s.maxLen > 0 {
		code = code[:s.maxLen]
	}


	 redirect = &Redirect{
	 	Code: code,
	 	LongUrl: url,
	 	ShortUrl: strings.TrimRight(s.shortUrl,"/")+"/"+code,
	 	CreatedAt: now,
	 }

	store, err := s.repository.Store(redirect)
	if err != nil {
		return
	}
	if store!=nil {
		chainInfo := &ChainInfo{}
		json.Unmarshal(store,chainInfo)
		redirect.ChainInfo = chainInfo
	}

	//redirect.LongUrl = strings.TrimRight(s.shortUrl,"/")+"/"+code
	return redirect,nil
}

func New(middleware []Middleware,repository Repository,logger log.Logger,shortUrl string,maxLen int) Service {
	var svc = NewService(logger,repository,shortUrl,maxLen)
	for _,mid := range middleware{
		svc = mid(svc)
	}
	return svc
}

func NewService(logger log.Logger, repository Repository, shortUrl string, maxLength int) Service {
	if maxLength > 9 {
		maxLength = 9
	}
	return &service{repository: repository, shortUrl: shortUrl, logger: logger, maxLen: maxLength}
}

