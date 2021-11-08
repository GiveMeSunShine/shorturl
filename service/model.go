/**
 * @Author : ysh
 * @Description :
 * @File : model
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 上午10:30
 */

package service

import "time"

type Redirect struct {
	Code      string    `json:"code"`
	LongUrl       string `json:"long_url"`
	ShortUrl  string `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}
