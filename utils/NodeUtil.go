/**
 * @Author : ysh
 * @Description : 长安链工具类
 * @File : NodeUtil
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/16 下午4:31
 */

package utils

import sdk "chainmaker.org/chainmaker/sdk-go/v2"

func CreateChainClientWithSDKConf(sdkConfPath string) (*sdk.ChainClient, error) {
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
	)
	if err != nil {
		return nil, err
	}

	// Enable certificate compression
	err = cc.EnableCertHash()
	if err != nil {
		return nil, err
	}
	return cc, nil
}
