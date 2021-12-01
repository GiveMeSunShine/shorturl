/**
 * @Author : ysh
 * @Description :
 * @File : localRepository
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/12/1 下午5:08
 */

package local

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"shorturl/service"
)

type localRepository struct {
	path string
}

func (l *localRepository) Upload(file string, fileName string) (hash string, err error) {
	fileBytes := bytes.NewBufferString(file).Bytes()
	sum := md5.Sum(fileBytes)
	hash = fmt.Sprintf("%x", sum)

	f, err := os.OpenFile(l.path+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, bytes.NewBufferString(file))

	return hash, nil
}

func (l *localRepository) Download(hash string) string {
	panic("implement me")
}

func NewLocalRepository(localPath string) (service.FileRepository, error) {
	log.Println("=========== Create Local File Repository ==============")
	return &localRepository{path: localPath}, nil
}
