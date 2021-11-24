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
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"log"
	"shorturl/service"
	"time"
)

type mysqlRepository struct {
	client   *xorm.Engine
}

func (m *mysqlRepository) Exists(has string) (exists bool, err error) {
	return
}

func NewMySQLRepository(drive string, hosts, port ,username, password, database string,idleConn int,maxConn int) (service.Repository, error) {
	mysqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8\n",username,password, hosts,port, "mysql")
	log.Println("conn mysql url :",mysqlUrl)
	mySQL, err := xorm.NewMySQL(drive, mysqlUrl)
	if err != nil {
		log.Fatal("xorm connect mysql Err : ",err)
		return nil, err
	}
	mySQL.Ping()

	CreateDatabase(mySQL)
	if err = mySQL.Close(); err != nil {
		log.Fatal("DB.close err : ", err)
		panic(err)
	}
	mysqlUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8\n",username,password, hosts,port, database)
	log.Println("conn mysql url :",mysqlUrl)
	mySQL, err = xorm.NewMySQL(drive, mysqlUrl)
	if err != nil {
		log.Fatal("xorm connect mysql Err : ",err)
		return nil, err
	}
	mySQL.Ping()
	mySQL.SetMaxOpenConns(maxConn)
	mySQL.SetConnMaxLifetime(100*time.Second)
	mySQL.SetMaxIdleConns(idleConn)

	CreateTable(mySQL)
	return &mysqlRepository{client: mySQL}, nil
}

func CreateDatabase(engine *xorm.Engine)  {
	log.Println("=====> start create database shorturl ")
	exec, dbErr := engine.DB().Exec("create database if not exists shorturl DEFAULT CHARSET utf8 COLLATE utf8_general_ci;")
	if dbErr != nil {
		log.Fatal("create Database Err : ",dbErr)
		return
	}
	id, dbErr := exec.LastInsertId()
	if dbErr != nil {
		log.Fatal("create Database LastInsertId Err : ",dbErr)
		return
	}
	log.Println("<===== create database result : ",id)
}

func CreateTable(engine *xorm.Engine) {
	log.Println("=====> start create url_map table")
	tableSql := "CREATE TABLE if not exists `url_map`( "+
		"`id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键', "+
		"`short_url` VARCHAR(64) NOT NULL COMMENT '短链URL', "+
		"`long_url` VARCHAR(768) NOT NULL COMMENT '长链URL', "+
		"`compression_code` VARCHAR(16) NOT NULL COMMENT '压缩码', "+
		"`description` VARCHAR(256) COMMENT '描述', "+
		"`url_status` TINYINT NOT NULL DEFAULT 1 COMMENT 'URL状态,1:正常,2:已失效', "+
		"`create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', "+
		"`edit_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间', "+
		"`creator` VARCHAR(32) NOT NULL DEFAULT 'admin' COMMENT '创建者', "+
		"`editor` VARCHAR(32) NOT NULL DEFAULT 'admin' COMMENT '更新者', "+
		"`deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '软删除标识', "+
		"`version` BIGINT NOT NULL DEFAULT 1 COMMENT '版本号', "+
		"UNIQUE uniq_compression_code (`compression_code`), "+
		"INDEX idx_short_url (`short_url`) "+
		") COMMENT 'URL映射表' ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin; "
	_, dbErr := engine.DB().Exec(tableSql)
	if dbErr !=nil {
		log.Fatal("create url_map table Err : ",dbErr)
	}

	log.Println("<===== create url_map table end ")
}



func (m *mysqlRepository) Find(code string) (redirect *service.Redirect, err error) {
	log.Println("[MySQL] Find : ",code)
	findSql := fmt.Sprintf("select short_url,long_url,compression_code from url_map where compression_code='%s' and  deleted=0 and url_status=1",
		code)
	row := m.client.DB().QueryRow(findSql)
	redirect = new(service.Redirect)
	err = row.Scan(&redirect.ShortUrl, &redirect.LongUrl, &redirect.Code)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}

func (m *mysqlRepository) Store(redirect *service.Redirect) ([]byte,error) {
	insertSql := fmt.Sprintf("insert into url_map (short_url,long_url,compression_code) values('%s','%s','%s') ",
		redirect.ShortUrl, redirect.LongUrl, redirect.Code)
	log.Println("insertSql : ",insertSql)
	exec, err := m.client.DB().Exec(insertSql)
	if err != nil {
		return nil,err
	}
	id, err := exec.LastInsertId()
	if err != nil {
		return nil,err
	}
	fmt.Println("[MySQL] Store success : ",id)
	return nil,nil
}
