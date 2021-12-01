/**
 * @Author : ysh
 * @Description :
 * @File : model
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 上午10:30
 */

package service

import (
	"mime/multipart"
	"time"
)

type Redirect struct {
	Code      string    `json:"code"`
	OrgLink   string    `json:"org_link"`
	ShortUrl  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
	ChainInfo interface{}
}

type DBStore struct {
	Code      string    `json:"code"`
	OrgLink   string    `json:"org_link"`
	ShortUrl  string    `json:"short_url"`
	Type      string    `json:"type"`
	FileName  string    `json:"file_name"`
	CreatedAt time.Time `json:"created_at"`
}

type ChainInfo struct {
	Txid        string `json:"txid"`
	BlockHeight uint64 `json:"block_height"`
}

type PostInfo struct {
	LongUrl    string
	Type       string
	File       multipart.File
	FileHeader multipart.FileHeader
}
