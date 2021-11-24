/**
 * @Author : ysh
 * @Description :
 * @File : server
 * @Software: GoLand
 * @Version: 1.0.0
 * @Time : 2021/11/5 下午2:28
 */

package server

import (
	"fmt"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/oklog/oklog/pkg/group"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"os"
	"os/signal"
	"shorturl/configMgr"
	"shorturl/endpoint"
	"shorturl/logging"
	chainmaker "shorturl/repository/chainmaker"
	mysqlrepo "shorturl/repository/mysql"
	redisrepo "shorturl/repository/redis"
	"shorturl/service"
	httptransport "shorturl/transport/http"
	"strconv"
	"syscall"
	"time"
)

var logger log.Logger

var(
	/*fs            = flag.NewFlagSet("hello", flag.ExitOnError)
	httpAddr      = fs.String("http-addr", ":8080", "HTTP listen address")
	devCors       = fs.String("dev-cors", "false", "is develop")
	dbDrive       = fs.String("db-drive", "mysql", "db drive type, default: mysql")
	redisDrive    = fs.String("redis-drive", "single", "redis drive: single or cluster")
	redisHosts    = fs.String("redis-hosts", "localhost:6379", "redis hosts, many ';' split")
	redisPassword = fs.String("redis-password", "", "redis password")
	redisDB       = fs.String("redis-db", "0", "redis db")
	mysqlHosts    = fs.String("mysql-hosts", "127.0.0.1:3306", "mysql hosts")
	mysqlPassword = fs.String("mysql-password", "root", "mysql password")
	mysqlDB       = fs.String("mysql-db", "shorturl", "mysql db")
	shortUrl      = fs.String("shortid-url", "http://localhost:8080/", "shortid url")
	logPath       = fs.String("log-path", "", "logging file path.")
	logLevel      = fs.String("log-level", "all", "logging level.")
	rateBucketNum = fs.Int("rate-bucket", 10, "rate bucket num")
	maxLength     = fs.Int("max-length", -1, "code length")*/
	err           error
)

type Server struct {
	conf *configMgr.ShortUrlConf
}

func Run() {
	/*if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}
	redisDrive = envString("REDIS_DRIVE", redisDrive)
	redisHosts = envString("REDIS_HOSTS", redisHosts)
	redisPassword = envString("REDIS_PASSWORD", redisPassword)
	redisDB = envString("REDIS_DB", redisDB)
	mysqlHosts    = envString("MYSQL_HOSTS", mysqlHosts)
	mysqlPassword = envString("MYSQL_PASSWORD", mysqlPassword)
	mysqlDB       = envString("MYSQL_DB", mysqlDB)
	dbDrive = envString("DB_DRIVE", dbDrive)
	shortUrl = envString("SHORT_URL", shortUrl)
	logPath = envString("LOG_PATH", logPath)
	logLevel = envString("LOG_LEVEL", logLevel)
	devCors = envString("DEV_CORS", devCors)
	rateBucketNum = envInt("RATE_BUCKET", rateBucketNum)
	maxLength = envInt("MAX_LENGTH", maxLength)*/
	conf, err := configMgr.NewShortUrlConf()
	if err != nil {
		logger.Log("read conf Err :",err)
		panic(err)
		return
	}
	viper := conf.GetShortConfigViper()
	server:= &Server{
		conf: conf,
	}
	dBType := viper.GetString("db.dBType")
	logPath := viper.GetString("log.path")
	logLevel := viper.GetString("log.level")

	logger = logging.SetLogging(logger, &logPath, &logLevel)

	var repo service.Repository
	switch dBType {
	case "mysql":
		hosts := viper.GetString("db.mysql.ip")
		port := viper.GetString("db.mysql.port")
		username := viper.GetString("db.mysql.username")
		passwd := viper.GetString("db.mysql.passwd")
		database := viper.GetString("db.mysql.libName")
		maxConn := viper.GetInt("db.mysql.maxConn")
		idleConn := viper.GetInt("db.mysql.idleConn")
		repo, err = mysqlrepo.NewMySQLRepository(dBType,hosts,port,username,passwd,database,idleConn,maxConn)
		if err != nil {
			_ = level.Error(logger).Log("connect", "db", "err", err.Error())
			return
		}
	case "redis":
		redisDB := viper.GetString("db.redis.libName")
		redisDrive := viper.GetString("db.redis.drive")
		redisHosts := viper.GetString("db.redis.hosts")
		redisPassword := viper.GetString("db.redis.passwd")
		db, _ := strconv.Atoi(redisDB)
		repo, err = redisrepo.NewRedisRepository(redisrepo.RedisDrive(redisDrive), redisHosts, redisPassword, "shorter", db)
		if err != nil {
			_ = level.Error(logger).Log("connect", "db", "err", err.Error())
			return
		}
	case "chainmaker":
		conteactName := viper.GetString("db.chainmaker.conteactName")
		configPath := viper.GetString("db.chainmaker.configpath")
		repo, err =chainmaker.NewMakerRepository(conteactName,configPath)
		if err != nil {
			_ = level.Error(logger).Log("connect", "chainmaker", "err", err.Error())
			return
		}
	}
	shortUrl := viper.GetString("short.defaultUrl")
	maxLength := viper.GetInt("short.maxLength")

	svc := service.New(getServiceMiddleware(logger),repo,logger,shortUrl,maxLength)
	eps := endpoint.New(svc, server.getEndpointMiddleware(logger))
	g := server.createService(eps)
	initCancelInterrupt(g)
	_ = logger.Log("exit", g.Run())
}

func (server *Server)createService(endpoints endpoint.Endpoints) (g *group.Group) {
	g = &group.Group{}
	server.initHttpHandler(endpoints, g)
	return g
}

func (server *Server)initHttpHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultHttpOptions(logger)

	httpHandler := httptransport.NewHttpHandler(endpoints,options)
	httpAddr  := server.conf.Viper.GetString("short.http.address")
	devCors := server.conf.Viper.GetString("develop")
	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		_ = level.Error(logger).Log("transport", "HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		_ = level.Debug(logger).Log("transport", "HTTP", "addr", httpAddr)
		headers := make(map[string]string)
		if isDev, _ := strconv.ParseBool(devCors); isDev {
			headers = map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS,PUT,DELETE",
				"Access-Control-Allow-Headers": "Origin,Content-Type,mode,Authorization,x-requested-with,Access-Control-Allow-Origin,Access-Control-Allow-Credentials",
			}
		}
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))
		return http.Serve(httpListener, accessControl(httpHandler, logger, headers))
	}, func(error) {
		_ = httpListener.Close()
	})
}

func defaultHttpOptions(logger log.Logger) map[string][]kithttp.ServerOption {
	options := map[string][]kithttp.ServerOption{"Get": {
		kithttp.ServerErrorEncoder(httptransport.ErrorRedirect),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	},
		"Post": {
			kithttp.ServerErrorEncoder(httptransport.ErrorEncoder),
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerBefore(kithttp.PopulateRequestContext),
		}}

	return options
}
func accessControl(h http.Handler, logger log.Logger, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})
}

func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	mw = addDefaultServiceMiddleware(logger, mw)
	return
}
func (server *Server)getEndpointMiddleware(logger log.Logger) (mw map[string][]kitendpoint.Middleware) {
	mw = map[string][]kitendpoint.Middleware{}
	mw = server.addDefaultEndpointMiddleware(logger, mw)

	return
}

func envString(env string, fallback *string) *string {
	e := os.Getenv(env)
	if e == "" {
		_ = os.Setenv(env, *fallback)
		return fallback
	}
	return &e
}

func envInt(env string, fallback *int) *int {
	e := os.Getenv(env)
	if e == "" {
		_ = os.Setenv(env, strconv.Itoa(*fallback))
		return fallback
	}
	num, _ := strconv.Atoi(e)
	return &num
}

func addDefaultServiceMiddleware(logger log.Logger, mw []service.Middleware) []service.Middleware {
	mw = append(mw, service.LoggingMiddleware(logger))
	return mw
}

func (server *Server)addDefaultEndpointMiddleware(logger log.Logger, mw map[string][]kitendpoint.Middleware) map[string][]kitendpoint.Middleware {
	rateBucketNum := server.conf.Viper.GetInt("short.rateBucketNum")
	mw["Post"] = append(mw["Post"],
		endpoint.LoggingMiddleware(logger),
		endpoint.NewTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum)),
	)

	mw["Get"] = append(mw["Get"],
		endpoint.LoggingMiddleware(logger),
		endpoint.NewTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), rateBucketNum*100)),
	)

	return mw
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}