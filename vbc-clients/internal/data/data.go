package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"vbc/configs"
	"vbc/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewCommonRepo)

// Data .
type Data struct {
	Db          *gorm.DB
	RedisClient *redis.Client
	PGSQL       *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)
	cleanup := func() {
		logHelper.Info("closing the data resources")
	}
	d := &Data{}
	if configs.AppName == configs.App_vbc {
		d.Db = mysqlInit(c)
		d.RedisClient = redisInit(c)
	} else {
		d.Db = mysqlInit(c)
		d.RedisClient = redisInit(c)
	}
	if configs.IsDev() {
		//s, err := postgresSqlInit(c)
		//if err != nil {
		//	logHelper.Error(err)
		//}
		//d.PGSQL = s
	}

	return d, cleanup, nil
}

//
//func createTLSConf(conf *conf.Data) tls.Config {
//
//	rootCertPool := x509.NewCertPool()
//	pem, err := os.ReadFile(conf.Database.GetCertPem())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
//		log.Fatal("Failed to append PEM.")
//	}
//	//clientCert := make([]tls.Certificate, 0, 1)
//	//
//	//certs, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//clientCert = append(clientCert, certs)
//
//	return tls.Config{
//		RootCAs: rootCertPool,
//		//Certificates:       clientCert,
//		InsecureSkipVerify: true, // needed for self signed certs
//	}
//}

func postgresSqlInit(conf *conf.Data) (*gorm.DB, error) {
	dsn := "host=postgresqleu2.postgres.database.azure.com user=pvbcuser password=Pz4kHsk8mXz2pDEc dbname=vbcvector port=5432 TimeZone=UTC"
	fmt.Println("postgresSqlInit")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func mysqlInit(conf *conf.Data) *gorm.DB {

	var db *gorm.DB
	var err error

	driver := "mysql"
	dsn := configs.EnvMySQLDSN()
	if configs.IsProd() {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DriverName: driver,
			DSN:        dsn,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // Silent: 不出现record not found
		})
	} else {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DriverName: driver,
			DSN:        dsn,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		panic(err)
	}
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(30)
	return db
}

func redisInit(conf *conf.Data) *redis.Client {
	//fmt.Println("redisInit:", conf.RedisUrl)
	dsn := configs.EnvRedisDSN()
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		panic(err)
	}
	opt.MaxActiveConns = 5
	opt.PoolSize = 5
	opt.DB = 0
	client := redis.NewClient(opt)
	ping := client.Ping(context.TODO())
	if ping.Err() != nil {
		if ping.Err() != redis.Nil {
			panic(ping.Err())
		}
	}
	return client
}
