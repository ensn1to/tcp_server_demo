
# TCP Server Demo

A simple tcp server demo which implements a application-level protocal based on tcp protocal.


## Protocaol

### Application level protocaol

#### Frame
```
totalLength(4)|data(n)

| segment    | type   | size    | remark                  |
| ---------- | ------ | ------- | ----------------------- |
| `dataSize` | int32 | 4       | the size of `data` only  |
| `data`     | []byte | dynamic |                         |
```

#### Packet
```
totalLength(4)|data(n)

| segment    | type   | size    | remark                  |
| ---------- | ------ | ------- | ----------------------- |
| `commandId`| byte   | 1       | message type            |
| `id`       | int64  | 8       | message id              |
| `payload`  | []byte | dynamic |                         |
```


### How to use it 

```
$ make
$ ./server
$ ./client
```