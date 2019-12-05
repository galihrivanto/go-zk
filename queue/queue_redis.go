package queue

import (
	"bufio"
	"bytes"
	"log"
	"strings"
	"time"

	driver "github.com/gomodule/redigo/redis"
)

// RedisOption define redis configuration for
// redis queue
type RedisOption struct {
	// host of redis
	Host string

	// credential to access redis
	Password string

	// DB index
	DB int

	// maximum maintained idle connection
	MaxIdle int

	// max active connection
	MaxActive int

	// maximum time connection being before closed
	IdleTimeout time.Duration

	// wait for new connection
	Wait bool
}

// RedisQueue implement Queuer using redis
type RedisQueue struct {
	pool *driver.Pool
}

// Push implements Queuer.Push
func (q *RedisQueue) Push(queueName string, item string) error {
	conn := q.pool.Get()
	defer conn.Close()

	return conn.Send("LPUSH", queueName, item)
}

// Pop implements Queuer.Pop
func (q *RedisQueue) Pop(queueName string) (string, error) {
	conn := q.pool.Get()
	defer conn.Close()

	return driver.String(conn.Do("RPOP", queueName))
}

// Len implements Queuer.Len
func (q *RedisQueue) Len(queueName string) int {
	conn := q.pool.Get()
	defer conn.Close()

	v, err := driver.Int(conn.Do("LLEN", queueName))
	if err != nil {
		return 0
	}

	return v
}

// NewRedisQueue init new queue which backed by redis
func NewRedisQueue(option *RedisOption) Queuer {
	if option == nil {
		option = &RedisOption{}
	}

	if option.Host == "" {
		option.Host = "localhost:6379"
	}

	// init pool
	pool, err := open(option)
	if err != nil {
		panic(err)
	}

	return &RedisQueue{pool: pool}
}

// open connection
func open(option *RedisOption) (*driver.Pool, error) {

	pool := &driver.Pool{
		MaxIdle:     option.MaxIdle,
		MaxActive:   option.MaxActive,
		IdleTimeout: option.IdleTimeout,
		Wait:        option.Wait,
		Dial: func() (driver.Conn, error) {
			c, err := driver.Dial("tcp", option.Host)
			if err != nil {
				return nil, err
			}

			if option.Password != "" {
				if _, err := c.Do("AUTH", option.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c driver.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	c := pool.Get()
	defer c.Close()

	v, err := driver.String(c.Do("INFO"))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(bytes.NewBufferString(v))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "redis_version") {
			log.Printf("\t== %s\n", line)
			break
		}
	}

	return pool, nil
}
