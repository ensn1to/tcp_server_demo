
# TCP Server Demo

A simple tcp server demo which implements a application-level protocal based on tcp protocal.


## Protocaol

### application level protocaol

```
Frame
totalLength(4)|data(n)

| segment    | type   | size    | remark                  |
| ---------- | ------ | ------- | ----------------------- |
| `dataSize` | int32 | 4       | the size of `data` only |                       |
| `data`     | []byte | dynamic |                         |
```