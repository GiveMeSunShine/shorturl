/**
 * @Author : ysh
 * @Description :
 * @File : server
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 下午2:28
 */

package redis

import (
	"fmt"
	"github.com/pkg/errors"
	"shorturl/service"
	"time"
)

type redisRepository struct {
	client RedisInterface
}

func (m *redisRepository) Exists(has string) (exists bool, err error) {
	return
}

func NewRedisRepository(drive RedisDrive, hosts, password, prefix string, database int) (service.Repository, error) {
	rdsClient := NewRedisClient(drive, hosts, password, prefix, database)

	return &redisRepository{client: rdsClient}, nil
}

func (m *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (m *redisRepository) Find(code string) (redirect *service.Redirect, err error) {
	data, err := m.client.HGetAll(m.generateKey(code))
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(service.DataNotFoundErr, "repository.Redirect.Find")
	}

	now, err := time.Parse("2006-01-02 15:04:05", data["created_at"])
	if err != nil {
		return
	}

	return &service.Redirect{
		Code:      data["code"],
		LongUrl:       data["long_url"],
		CreatedAt: now.In(time.Local),
		ShortUrl: data["short_url"],
	}, nil
}

func (m *redisRepository) Store(redirect *service.Redirect) ([]byte,error) {
	data := map[string]interface{}{
		"code":       redirect.Code,
		"long_url":        redirect.LongUrl,
		"created_at": redirect.CreatedAt.Format("2006-01-02 15:04:05"),
		"short_url": redirect.ShortUrl,
	}

	err := m.client.HMSet(m.generateKey(redirect.Code), data)
	if err != nil {
		return nil,errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil,nil
}
