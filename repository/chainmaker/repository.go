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
	ShortUrl string `json:"shortUrl"`
	LongUrl string `json:"longUrl"`
	Code string `json:"code"`
	Description string `json:"description"`
	Creator string `json:"creator"`
	Version string `json:"version"`
	Time     int32 `json:"time"`
}

type makerRepository struct {
	client   *sdk.ChainClient
	contractName string
}

func (m *makerRepository) Exists(has string) (exists bool, err error) {
	return
}

func NewMakerRepository(contractName string,configPath string) (service.Repository, error) {
	log.Println("=========== Create ChainClient ==============")
	client, err := utils.CreateChainClientWithSDKConf(configPath)
	if err!=nil {
		log.Fatalln("Create Chain Client ERR:",err)
	}
	height, err := client.GetCurrentBlockHeight()
	if err != nil {
		log.Fatalln("Create Chain Client ERR:",err)
	}
	log.Println(" ===> Current Block Height :",height )
	return &makerRepository{client: client,contractName: contractName}, nil
}


func (m *makerRepository) Find(code string) (redirect *service.Redirect, err error) {
	log.Println("[chainMaker] Find : ",code)
	kvs := []*common.KeyValuePair{
		{
			Key:   "code",
			Value: []byte(code),
		},
	}
	query, err := userContractClaimQuery(m.client, m.contractName, "find_by_code", kvs)
	if err != nil {
		log.Fatalln("[chainMaker]  userContractClaimQuery Err :",err)
	}
	redirect = &service.Redirect{
		Code: query.Code,
		ShortUrl: query.ShortUrl,
		LongUrl: query.LongUrl,
	}
	return redirect, nil
}

func (m *makerRepository) Store(redirect *service.Redirect) ([]byte,error) {

	invoke, err := userContractClaimInvoke(m.client, m.contractName,redirect, "save", true)
	if err != nil {
		log.Fatalln("[chainMaker] Store ERR :",err)
	}
	log.Println("[chainMaker] Store success , Txid : ",string(invoke))
	return invoke,nil
}

func userContractClaimInvoke(client *sdk.ChainClient, contractName string, redirect *service.Redirect, method string, withSyncResult bool) ([]byte, error) {
	curTime := strconv.FormatInt(redirect.CreatedAt.Unix(), 10)
	kvs := []*common.KeyValuePair{
		{
			Key:   "time",
			Value: []byte(curTime),
		},
		{
			Key:   "short_url",
			Value: []byte(redirect.ShortUrl),
		},
		{
			Key:   "long_url",
			Value: []byte(redirect.LongUrl),
		},
		{
			Key:   "code",
			Value: []byte(redirect.Code),
		},{
			Key:   "description",
			Value: []byte("测试数据"),
		},{
			Key:   "creator",
			Value: []byte("admin"),
		},{
			Key:   "version",
			Value: []byte("1.0.0"),
		},
	}

	result, err := invokeUserContract(client, contractName, method, "", kvs, withSyncResult)
	if err != nil {
		log.Fatalln("[chainMaker] invokeUserContract Err : ",err)
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
		Txid: resp.GetTxId(),
		BlockHeight: info.BlockHeight,
	}

	/*id, err := client.GetBlockByTxId(resp.GetTxId(), true)
	if err != nil {
		log.Fatalln("GetBlockByTxId Err : ",err)
	}
	log.Println("GetBlockByTxId : ",id)

	byTxId, err := client.GetTxByTxId(resp.GetTxId())
	if err != nil {
		log.Fatalln("GetTxByTxId Err : ",err)
	}
	log.Println("GetTxByTxId : ",byTxId)*/

	marshal, _ := json.Marshal(chainInfo)

	return marshal,nil
}

func userContractClaimQuery(client *sdk.ChainClient,contractName string,method string, kvs []*common.KeyValuePair) (short *Short, err error){
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
			return short,err
		}
		return short,nil
	}else {
		log.Fatalln("QUERY claim contract err")
	}
	return short,nil
}
