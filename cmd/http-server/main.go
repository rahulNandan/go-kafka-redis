package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gitlab.com/melwyn95/go-redis-kafka-demo/pkg/kafka"
	"gitlab.com/melwyn95/go-redis-kafka-demo/pkg/redis"

	"github.com/gorilla/mux"
)

const MILLION int = 1000000

type ProduceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (sc *ServerConfig) produceHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	b, _ := json.Marshal(map[string]interface{}{
		"type":   "number",
		"number": sc.rng.Intn(MILLION),
	})

	err := sc.kafka.Produce(b)

	var status, message string
	if err != nil {
		status = "failure"
		message = fmt.Sprintf("unable to produced event on kafka topic %s", *sc.kafka.Topic)
	} else {
		status = "success"
		message = fmt.Sprintf("event produced on kafka topic %s", *sc.kafka.Topic)
	}

	json.NewEncoder(&buf).Encode(&ProduceResponse{
		Status:  status,
		Message: message,
	})

	io.WriteString(w, buf.String())
}

type PingResponse struct {
	NoOfRequests int `json:"requests"`
	Even         int `json:"even"`
	Odd          int `json:"odd"`
}

func (sc *ServerConfig) pingHandler(w http.ResponseWriter, r *http.Request) {
	var evens, odds, requests int
	var buf bytes.Buffer

	c, _ := sc.redis.Read([]string{"even", "odd"})
	for i := range c {
		if c[i] == nil {
			break
		}
		if i == 0 {
			evens, _ = strconv.Atoi(c[i].(string))
		} else {
			odds, _ = strconv.Atoi(c[i].(string))
		}
	}
	requests = odds + evens

	json.NewEncoder(&buf).Encode(&PingResponse{
		NoOfRequests: requests,
		Even:         evens,
		Odd:          odds,
	})

	io.WriteString(w, buf.String())
}

type ServerConfig struct {
	redis *redis.Redis
	kafka *kafka.KafkaProducer
	rng   *rand.Rand
}

type ConsumerConfig struct {
	redis *redis.Redis
	kafka *kafka.KafkaConsumer
}

func gracefullyShutdownServer(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("could not gracefully shutdown http-server: %v\n", err)
	}

	close(done)
}

func main() {
	// redis config vars
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	// kafka config vars
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	consumerGroupID := os.Getenv("KAFKA_CONSUMER_GROUP")
	// http-server config vars
	serverPort := os.Getenv("HTTP_SERVER_PORT")
	serverTimeoutStr := os.Getenv("HTTP_SERVER_TIMEOUT")
	serverTimeout, _ := strconv.Atoi(serverTimeoutStr)

	quitServer := make(chan os.Signal, 1)
	quitConsumner := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(quitServer, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quitConsumner, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	consumerRedis, _ := redis.New("consumer", redisHost, redisPort, redisPassword)
	kafkaConsumer, _ := kafka.NewConsumer(kafkaHost, kafkaPort, kafkaTopic, consumerGroupID)
	// Spawn a go routine and listen to kafka topic
	go kafka.SpawnConsumer(kafkaConsumer, consumerRedis, quitConsumner)

	producerRedis, _ := redis.New("producer", redisHost, redisPort, redisPassword)
	kafkaProducer, _ := kafka.NewProducer(kafkaHost, kafkaPort, kafkaTopic)
	sc := ServerConfig{
		redis: producerRedis,
		kafka: kafkaProducer,
		rng:   rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
	}
	router := mux.NewRouter()

	// Get number of odd & even numbers produced on kafka
	router.HandleFunc("/", sc.pingHandler)

	// Produces an event on kafka with a number
	router.HandleFunc("/produce", sc.produceHandler)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", serverPort),
		Handler: http.TimeoutHandler(router, time.Duration(serverTimeout)*time.Second, "Server Timeout"),
	}

	go gracefullyShutdownServer(&srv, quitServer, done)

	fmt.Printf("http-server ready listening on port %s \n", serverPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("failed to listen:%v", err)
	}

	<-done
	fmt.Println("http-server exited gracefully...")

	// close redis connection pools
	producerRedis.Close()
	consumerRedis.Close()
	// close kafka producer
	kafkaProducer.Close()
}
