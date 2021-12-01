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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/log"
	"io"
	"shorturl/shortid"
	"strings"
	"time"
)

var (
	DataNotFoundErr = errors.New("Redirect Data Not Found")
	DataInvalidErr  = errors.New("Redirect Data Invalid")
)

type Service interface {
	Get(ctx context.Context, code string) (redirect *Redirect, err error)
	Post(ctx context.Context, parm PostInfo) (redirect *Redirect, err error)
}

type service struct {
	repository     Repository
	fileRepository FileRepository
	logger         log.Logger
	shortUrl       string
	maxLen         int
}

func (s service) Get(ctx context.Context, code string) (redirect *Redirect, err error) {
	find, err := s.repository.Find(code)
	if err != nil {
		//TODO
		return nil, err
	}
	redirect = &Redirect{
		OrgLink: find.OrgLink,
	}
	return redirect, nil
}

func (s service) Post(ctx context.Context, parm PostInfo) (redirect *Redirect, err error) {
	now := time.Now()
	now = now.In(time.Local)
	var code string
	code = shortid.ShortIdFromCode()
	if s.maxLen > 0 {
		code = code[:s.maxLen]
	}
	shortUrl := strings.TrimRight(s.shortUrl, "/") + "/" + code
	dbStore := &DBStore{
		Code:      code,
		OrgLink:   parm.LongUrl,
		ShortUrl:  shortUrl,
		Type:      parm.Type,
		CreatedAt: now,
	}
	if parm.Type == "1" {
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, parm.File); err != nil {
			return nil, err
		}
		hash, err := s.fileRepository.Upload(buf.String(), parm.FileHeader.Filename)
		if err != nil {
			// TODO
		}
		dbStore.OrgLink = hash
		dbStore.FileName = parm.FileHeader.Filename
	}
	redirect = &Redirect{
		Code:      code,
		OrgLink:   parm.LongUrl,
		ShortUrl:  shortUrl,
		CreatedAt: now,
	}
	store, err := s.repository.Store(dbStore)
	if err != nil {
		return
	}
	if store != nil {
		chainInfo := &ChainInfo{}
		json.Unmarshal(store, chainInfo)
		redirect.ChainInfo = chainInfo
	}

	return redirect, nil
}

func New(middleware []Middleware, repository Repository, fileRepository FileRepository, logger log.Logger, shortUrl string, maxLen int) Service {
	var svc = NewService(logger, repository, fileRepository, shortUrl, maxLen)
	for _, mid := range middleware {
		svc = mid(svc)
	}
	return svc
}

func NewService(logger log.Logger, repository Repository, fileRepository FileRepository, shortUrl string, maxLength int) Service {
	if maxLength > 9 {
		maxLength = 9
	}
	return &service{repository: repository, fileRepository: fileRepository, shortUrl: shortUrl, logger: logger, maxLen: maxLength}
}
