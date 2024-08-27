# Go Connection Pool Implementation and Benchmarking

This Go program demonstrates a custom implementation of a database connection pool, along with benchmarking functions to compare the performance of using pooled versus non-pooled connections.

## Overview

- **Connection Pooling**: Efficiently manages a fixed number of database connections, reducing the overhead of repeatedly establishing and closing connections.
- **Benchmarking**: Compares the performance of executing database operations with and without using a connection pool.

## Code Structure

### `conn` Struct

```go
type conn struct {
    db *sql.DB
}
```

- Represents a single database connection.

### `connPool` Struct

```go
type connPool struct {
    mute           *sync.Mutex
    connections    []*conn
    maxConnections int
    channel        chan interface{}
}
```

- **`mute`**: Mutex to ensure thread-safe access to the connection pool.
- **`connections`**: Slice of pointers to `conn`, representing available connections.
- **`maxConnections`**: Maximum number of connections allowed in the pool.
- **`channel`**: Buffered channel used to manage connection availability.

### `NewConnectionPool(maxConns int) (*connPool, error)`

```go
func NewConnectionPool(maxConns int) (*connPool, error)
```

- Initializes a new connection pool with a specified maximum number of connections.
- Pre-fills the pool with connections and returns a pointer to the pool.

### `Get() (*conn, error)`

```go
func (pool *connPool) Get() (*conn, error)
```

- Retrieves a connection from the pool.
- Locks the pool, removes a connection from the slice, and returns it.

### `Put(c *conn)`

```go
func (pool *connPool) Put(c *conn)
```

- Returns a connection back to the pool.
- Locks the pool, appends the connection to the slice, and signals availability via the channel.

### `Close()`

```go
func (pool *connPool) Close()
```

- Closes all connections in the pool and the channel.

### `benchmarkPool()`

```go
func benchmarkPool()
```

- Benchmarks database operations using the connection pool.
- Creates a pool with 10 connections.
- Spawns 1000 goroutines to execute a simple SQL query (`SELECT SLEEP(0.01);`) using pooled connections.
- Measures and logs the total execution time.

### `benchmarkNonPool()`

```go
func benchmarkNonPool()
```

- Benchmarks database operations without using a connection pool.
- Spawns 20 goroutines, each creating and closing its own database connection.
- Measures and logs the total execution time.

### `NewConn() *sql.DB`

```go
func NewConn() *sql.DB
```

- Establishes a new connection to the MySQL database.
- Returns a pointer to the `sql.DB` object.

### `main()`

```go
func main()
```

- The entry point of the program.
- By default, calls the `benchmarkPool()` function to demonstrate connection pooling.
- The `benchmarkNonPool()` function can be uncommented to compare the performance with non-pooled connections.

## How to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/Aman123at/connection-pool-go.git
   cd connection-pool-go
   ```

2. Update the MySQL connection string in the `NewConn()` function:
   ```go
   db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/demo")
   ```

3. Run the program:
   ```bash
   go run main.go
   ```

4. To benchmark without a connection pool, uncomment the `benchmarkNonPool()` function call in `main()` and comment out the `benchmarkPool()` function call.

## Conclusion

This program demonstrates the performance benefits of using a connection pool for database operations in Go. The custom connection pool implementation shows how to efficiently manage database connections, and the benchmarking functions provide a clear comparison between pooled and non-pooled connection strategies.
