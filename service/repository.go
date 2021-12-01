/**
 * @Author : ysh
 * @Description : 数据存统一接口
 * @File : repository
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 上午10:31
 */

package service

type Repository interface {
	Find(code string) (redirect *DBStore, err error)
	Store(redirect *DBStore) ([]byte, error)
	Exists(has string) (exists bool, err error)
}

type FileRepository interface {
	Upload(file string, fileName string) (hash string, err error)
	Download(hash string) string
}
