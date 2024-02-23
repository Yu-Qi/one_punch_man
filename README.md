# one_punch_man

Backend practice with golang

## Practice

### Graceful shutdown

- [x] 透過 channel 接收 os signal 並進行 graceful shutdown
- [x] 透過 context 控制 goroutine 的生命週期，包含關閉子 goroutine
- [x] 透過 context 控制 http server 的生命週期，包含關閉 server

### Message queue

- []支援 http request 的方式
- []tcp 的方式
- [x]不保證順序性
- [x]當時間到或者數量到時，進行批次處理
- [x] 透過 for select 進行非同步處理
