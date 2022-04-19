# SQLDB 

SQLDB is sql db pool connection, this module can open multiple connection with key of your connection database.

## Base knowledge

### DSN

SQL Connection configuration is using dsn http url format: 
1. Our standard http dsn configuration, [here](https://github.com/xo/dburl).
   
    [list db dsn](https://github.com/xo/dburl/blob/master/dburl.go#L51).

    e.g : `postgres://postgres:12345678@127.0.0.1:5432/postgres?sslmode=disable`
2. Additional query params passing to `database/sql`
   - MaxConnLifeTime key param `max-ttl-conn` value int in second, [source](https://golang.org/pkg/database/sql/#DB.SetConnMaxLifetime). 
   - MaxIdleConnections key param `max-idle-conn` value int, [source](https://golang.org/pkg/database/sql/#DB.SetMaxIdleConns).
   - MaxOpenConnections key param `max-conn` value int, [source](https://golang.org/pkg/database/sql/#DB.SetMaxOpenConns).
   
    e.g : `?max-ttl-conn=60&max-idle-conn=1&max-conn=10`
### Driver

 Sql database driver support, import your database driver (at least once). You can use anonymous import also.
 [list driver](https://github.com/xo/dburl/blob/master/dburl.go#L129)

e.g : `import _ "github.com/lib/pq"`

### Open

```go
    Open(&config)
```

Open is open pool connection by passing configuration.
Ideally open before an application use pool

### Close

```go
    Close(
        key,//key is optional
    )
```
Close is close all connection or close by specific key.
Ideally close after gracefully shutdown service / before application shutdown 

### DB

```go
    DB(
        key,
    )
```
DB is get pool sql db by key. You can use pool for common sql query.
No need close after call DB