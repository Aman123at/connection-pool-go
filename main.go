package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type conn struct {
	db *sql.DB
}

type connPool struct {
	mute           *sync.Mutex
	connections    []*conn
	maxConnections int
	channel        chan interface{}
}

func NewConnectionPool(maxConns int) (*connPool, error) {
	var m = sync.Mutex{}
	pool := &connPool{
		mute:           &m,
		connections:    make([]*conn, 0, maxConns),
		maxConnections: maxConns,
		channel:        make(chan interface{}, maxConns),
	}
	for i := 0; i < maxConns; i++ {
		pool.connections = append(pool.connections, &conn{NewConn()})
		pool.channel <- nil
	}
	return pool, nil
}

func (pool *connPool) Close() {
	close(pool.channel)
	for i := range pool.connections {
		pool.connections[i].db.Close()
	}
}

func (pool *connPool) Get() (*conn, error) {
	<-pool.channel
	pool.mute.Lock()
	c := pool.connections[0]
	pool.connections = pool.connections[1:]
	pool.mute.Unlock()
	return c, nil
}

func (pool *connPool) Put(c *conn) {
	pool.mute.Lock()
	pool.connections = append(pool.connections, c)
	pool.mute.Unlock()
	pool.channel <- nil
}

func benchmarkPool() {
	start := time.Now()
	pool, err := NewConnectionPool(10)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := pool.Get()
			if err != nil {
				log.Fatal(err)
			}
			_, execerr := conn.db.Exec("SELECT SLEEP(0.01);")
			if execerr != nil {
				log.Fatal(err.Error())
			}
			pool.Put(conn)
		}()
	}
	wg.Wait()
	fmt.Println("Benchmark Connection Pool : ", time.Since(start))
	pool.Close()
}

func benchmarkNonPool() {
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db := NewConn()
			defer db.Close()
			_, err := db.Exec("SELECT SLEEP(0.01);")
			if err != nil {
				log.Fatal(err.Error())
			}

		}()
	}
	wg.Wait()
	fmt.Printf("BenchMark Non pool connections : %v", time.Since(startTime))
}

func NewConn() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/demo")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	log.Println("Welcome to connection pool")
	// benchmarkNonPool()
	benchmarkPool()
}
