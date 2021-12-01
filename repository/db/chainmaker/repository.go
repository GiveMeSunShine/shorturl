/**
 * @Author : ysh
 * @Description :
 * @File : repository
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/17 下午4:12
 */

package chainmaker

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"encoding/json"
	"fmt"
	"shorturl/utils"
	"strconv"

	"log"
	"shorturl/service"
)

type Short struct {
	ShortUrl    string `json:"shortUrl"`
	LongUrl     string `json:"longUrl"`
	Code        string `json:"code"`
	Type        string `json:"type"`
	FileName    string `json:"file_name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Version     string `json:"version"`
	Time        int32  `json:"time"`
}

type makerRepository struct {
	client       *sdk.ChainClient
	contractName string
}

func (m *makerRepository) Exists(has string) (exists bool, err error) {
	return
}

func NewMakerRepository(contractName string, configPath string) (service.Repository, error) {
	log.Println("=========== Create ChainClient ==============")
	client, err := utils.CreateChainClientWithSDKConf(configPath)
	if err != nil {
		log.Fatalln("Create Chain Client ERR:", err)
	}
	height, err := client.GetCurrentBlockHeight()
	if err != nil {
		log.Fatalln("Create Chain Client ERR:", err)
	}
	log.Println(" ===> Current Block Height :", height)
	return &makerRepository{client: client, contractName: contractName}, nil
}

func (m *makerRepository) Find(code string) (dbStore *service.DBStore, err error) {
	log.Println("[chainMaker] Find : ", code)
	kvs := []*common.KeyValuePair{
		{
			Key:   "code",
			Value: []byte(code),
		},
	}
	query, err := userContractClaimQuery(m.client, m.contractName, "find_by_code", kvs)
	if err != nil {
		log.Fatalln("[chainMaker]  userContractClaimQuery Err :", err)
	}
	dbStore = &service.DBStore{
		Code:     query.Code,
		ShortUrl: query.ShortUrl,
		OrgLink:  query.LongUrl,
		Type:     query.Type,
		FileName: query.FileName,
	}
	return dbStore, nil
}

func (m *makerRepository) Store(dbStore *service.DBStore) ([]byte, error) {

	invoke, err := userContractClaimInvoke(m.client, m.contractName, dbStore, "save", true)
	if err != nil {
		log.Fatalln("[chainMaker] Store ERR :", err)
	}
	log.Println("[chainMaker] Store success , Txid : ", string(invoke))
	return invoke, nil
}

func userContractClaimInvoke(client *sdk.ChainClient, contractName string, dbStore *service.DBStore, method string, withSyncResult bool) ([]byte, error) {
	curTime := strconv.FormatInt(dbStore.CreatedAt.Unix(), 10)
	kvs := []*common.KeyValuePair{
		{
			Key:   "time",
			Value: []byte(curTime),
		},
		{
			Key:   "short_url",
			Value: []byte(dbStore.ShortUrl),
		},
		{
			Key:   "long_url",
			Value: []byte(dbStore.OrgLink),
		},
		{
			Key:   "code",
			Value: []byte(dbStore.Code),
		}, {
			Key:   "type",
			Value: []byte(dbStore.Type),
		}, {
			Key:   "file_name",
			Value: []byte(dbStore.FileName),
		}, {
			Key:   "description",
			Value: []byte(""),
		}, {
			Key:   "creator",
			Value: []byte("admin"),
		}, {
			Key:   "version",
			Value: []byte("1.0.0"),
		},
	}

	result, err := invokeUserContract(client, contractName, method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln("[chainMaker] invokeUserContract Err : ", err)
		return nil, err
	}
	return result, nil
}

func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, kvs []*common.KeyValuePair, withSyncResult bool) ([]byte, error) {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
	if err != nil {
		return nil, err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return nil, fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if !withSyncResult {
		log.Println("[chainMaker] invoke contract success, resp: \n", resp.Code, resp.Message, resp.ContractResult)
	} else {
		log.Println("[chainMaker] invoke contract success, resp: \n", resp.Code, resp.Message, resp.ContractResult)
	}
	info, _ := client.GetChainInfo()
	chainInfo := &service.ChainInfo{
		Txid:        resp.GetTxId(),
		BlockHeight: info.BlockHeight,
	}
	marshal, _ := json.Marshal(chainInfo)

	return marshal, nil
}

func userContractClaimQuery(client *sdk.ChainClient, contractName string, method string, kvs []*common.KeyValuePair) (short *Short, err error) {
	resp, err := client.QueryContract(contractName, method, kvs, -1)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("QUERY claim contract resp: %+v\n", resp)
	if resp.Code == 0 {
		result := resp.ContractResult.Result
		short = &Short{}
		err := json.Unmarshal(result, &short)
		if err != nil {
			return short, err
		}
		return short, nil
	} else {
		log.Fatalln("QUERY claim contract err")
	}
	return short, nil
}
