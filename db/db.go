package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"senita-api/models"

	"github.com/go-gorp/gorp"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" //import postgres
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB ...
// type DB struct {
// 	*sql.DB
// }

var DB *gorm.DB

var db *gorp.DbMap

// var DV *gorm.DB

// Init ...
func Init() {
	dbinfo := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbinfo,
		PreferSimpleProtocol: true, // Use simple protocol for better performance
	}), &gorm.Config{})

	// db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{

	// 	CreateBatchSize: 1000,
	// })
	sqlDB, err := db.DB()
	if err != nil {
		// return nil, err
	}

	sqlDB.SetMaxIdleConns(10)  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100) // Maximum number of open connections
	fmt.Println(err)
	DB = db

}

// ConnectDB ...
func ConnectDB(dataSourceName string) (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	//dbmap.TraceOn("[gorp]", log.New(os.Stdout, "golang-gin:", log.Lmicroseconds)) //Trace database requests
	return dbmap, nil
}

// GetDB ...
func GetDB() *gorp.DbMap {
	return db
}

// RedisClient ...
var RedisClient *redis.Client

// InitRedis ...
func InitRedis(selectDB ...int) error {
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	db := 0 // Default DB value
	if len(selectDB) > 0 {
		db = selectDB[0]
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       db,
		// DialTimeout:        10 * time.Second,
		// ReadTimeout:        30 * time.Second,
		// WriteTimeout:       30 * time.Second,
		// PoolSize:           10,
		// PoolTimeout:        30 * time.Second,
		// IdleTimeout:        500 * time.Millisecond,
		// IdleCheckFrequency: 500 * time.Millisecond,
		// TLSConfig: &tls.Config{
		//     InsecureSkipVerify: true,
		// },
	})

	// Test the connection to Redis by sending a PING command
	pinged, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	// 	Ping the Redis server
	fmt.Println(pinged)
	return nil
}

// GetRedis ...

func AutoMigrate() {
	dbinfo := fmt.Sprintf("port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbinfo,
		PreferSimpleProtocol: true, // Use simple protocol for better performance
	}), &gorm.Config{})

	// db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{

	// 	CreateBatchSize: 1000,
	// })
	sqlDB, err := db.DB()
	if err != nil {
		// return nil, err
	}

	sqlDB.SetMaxIdleConns(1000)  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(10000) // Maximum number of open connections
	fmt.Println(err)
	DB = db

	db.Debug().AutoMigrate(&models.User{})
	db.Debug().AutoMigrate(&models.Category{})
	db.Debug().AutoMigrate(&models.Article{})

}
