# Concurrent safe Counter API

A concurrent safe API counter.

More improvements:
* Add tests for WS connections
* Add tests for events
* Do a faster cache for ws connections.


## How to run:
```
make build
./main

or 

make up
```

## Test flow:
```
open in terminal #1:
wscat -c "ws://localhost:8080/v1/subscribe"

open in terminal #2:
wscat -c "ws://localhost:8080/v1/subscribe"

open in terminal #3:
curl -X GET  http://localhost:8080/v1/a1
curl -X POST http://localhost:8080/v1/a1/increment
curl -X POST http://localhost:8080/v1/a1/increment
curl -X GET  http://localhost:8080/v1/a1
curl -X POST http://localhost:8080/v1/a1/decrement
curl -X GET  http://localhost:8080/v1/a1
curl -X POST http://localhost:8080/v1/a1/reset
curl -X GET  http://localhost:8080/v1/a1
```

## API:

### Get Stats by key:
```
GET /v1/:key
```
Example:
```
curl -X GET  http://localhost:8080/v1/a1
```

### Increment:
```
GET /v1/:key/increment
```
Example:
```
curl -X POST http://localhost:8080/v1/a1/increment
```

### Decrement:
```
GET /v1/:key/decrement
```
Example:
```
curl -X POST http://localhost:8080/v1/a1/decrement
```

### Reset:
```
GET /v1/:key/reset
```
Example:
```
curl -X POST http://localhost:8080/v1/a1/reset
```

### WS stats:
```
ws: /v1/subscribe
```
Example:
```
npm install -g wscat

wscat -c "ws://localhost:8080/v1/subscribe"
```