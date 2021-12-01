/**
 * @Author : ysh
 * @Description :
 * @File : ipfsRepository
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/12/1 下午5:07
 */

package ipfs

import (
	"bytes"
	ipfshell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"log"
	"shorturl/service"
)

type ipfsRepository struct {
	client *ipfshell.Shell
}

func (ipfs *ipfsRepository) Upload(file string, fileName string) (hash string, err error) {
	log.Println("upload file to IPFS ......")
	hash, err = ipfs.client.Add(bytes.NewBufferString(file))
	if err != nil {
		log.Fatalln("add to IPFS ERR :", err)
		return "", nil
	}
	log.Println("add to IPFS success! hash => ", hash)
	return hash, nil
}

func (ipfs *ipfsRepository) Download(hash string) string {
	log.Println("find file from PFS : hash => ", hash)
	cat, err := ipfs.client.Cat(hash)
	if err != nil {
		log.Fatalln("find from IPFS Err : ", err)
		return ""
	}
	all, err := ioutil.ReadAll(cat)
	if err != nil {
		log.Fatalln("ioutil.ReadAll Err : ", err)
		return ""
	}
	return string(all)
}

func NewIpfsRepository(ipfsAdd string) (service.FileRepository, error) {
	log.Println("=========== Create IPFS Client ==============")
	client := ipfshell.NewShell(ipfsAdd)
	return &ipfsRepository{client: client}, nil
}
