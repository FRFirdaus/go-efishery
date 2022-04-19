# requestidmw 

requestidmw is middleware support inject request/response header with key "request-id" and value uuid. request-id will use for micro services log tracking.

[problem](https://efishery.slack.com/archives/G4XA6GP9R/p1613040548028000)

## Example

- [gin](https://bitbucket.org/efishery/go-efishery/src/master/server/middleware/requestidmw/example/gin/main.go)
- [net/http](https://bitbucket.org/efishery/go-efishery/src/master/server/middleware/requestidmw/example/nethttp/main.go)

## Base knowledge

### HttpMiddleware

HttpMiddleware is middleware handler inject request/response using net/http lib. 
e.g

Using raw net/http
```go
import "bitbucket.org/efishery/go-efishery/server/middleware/requestidmw"

http.Handle("/hello", requestidmw.HttpMiddleware(handler()))
```
Using `github.com/gorilla/mux`
```go

	r := mux.NewRouter()
	r.Handle("/", handler())
	r.Handle("/hello", handler())

	// inject global middleware
	http.Handle("/", requestidmw.HttpMiddleware(r))
```

### GinMiddleware

GinMiddleware is middleware handler inject request/response using net/http lib. 
```go
import "bitbucket.org/efishery/go-efishery/server/middleware/requestidmw"

router := gin.Default()
// inject global middleware
router.Use(requestidmw.GinMiddleware())

router.GET("/hello", handler)
router.GET("/hello/me", handler)

// if only need specific endpoint
// router.GET("/hello/word",requestidmw.GinMiddleware(), handler)
```



### GetRequestId

GetRequestId is get request-id from request.
Required enable middleware 

```go

import "bitbucket.org/efishery/go-efishery/server/middleware/requestidmw"

// r is http.Request
requestidmw.GetRequestId(r)
```


### GetGinRequestId

GetGinRequestId is gin helper to get request-id from request.
Required enable middleware 

```go

import "bitbucket.org/efishery/go-efishery/server/middleware/requestidmw"

// ctx is gin context
requestidmw.GetGinRequestId(ctx)
```
