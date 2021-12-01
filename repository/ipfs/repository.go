/**
 * @Author : ysh
 * @Description :
 * @File : repository
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/30 上午10:55
 */

package ipfs

import (
	"bytes"
	ipfshell "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"log"
)

type ipfsRepository struct {
	client   *ipfshell.Shell
}

func NewIpfsRepository(ipfsAdd string) (*ipfsRepository, error) {
	log.Println("=========== Create IPFS Client ==============")
	client := ipfshell.NewShell(ipfsAdd)
	return &ipfsRepository{client: client}, nil
}

func (ipfs ipfsRepository) UploadIPFS(str string) string {
	log.Println("upload file to IPFS ......")
	hash, err := ipfs.client.Add(bytes.NewBufferString(str))
	if err != nil {
		log.Fatalln("add to IPFS ERR :",err)
		return ""
	}
	log.Println("add to IPFS success! hash => ",hash)
	return hash
}


func (ipfs *ipfsRepository)FindFile(hash string) string {
	log.Println("find file from PFS : hash => ",hash)
	cat, err := ipfs.client.Cat(hash)
	if err != nil {
		log.Fatalln("find from IPFS Err : ",err)
		return ""
	}
	all, err := ioutil.ReadAll(cat)
	if err != nil {
		log.Fatalln("ioutil.ReadAll Err : ",err)
		return ""
	}
	return string(all)
}





