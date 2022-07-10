package Connections

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "Revelt"
	password = ""
	dbname   = "kaktus"
)

var PostgreConn *sql.DB
var RabbitChannel *amqp.Channel
var RedisConn *redis.Client

func redisConn() (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:55001",
		Password: "redispw",
		DB:       0,
	})
	return
}

func postgresConn() (db *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	return
}

func rabbitConn() (con *amqp.Connection, err error) {
	con, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	return
}

func InitiateConnection() (err error) {
	PostgreConn, err = postgresConn()
	if err != nil {
		log.Panicf("[ERROR] : %s", err)
		return
	}
	rcon, err := rabbitConn()
	if err != nil {
		log.Panicf("[ERROR] : %s", err)
		return
	}
	RabbitChannel, err = rcon.Channel()
	if err != nil {
		log.Panicf("[ERROR] : %s", err)
		return
	}
	RedisConn = redisConn()
	return
}
