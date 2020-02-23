# PingCAP txn-store 作业

## 题目描述

作业:

KV store 里有 1000 个 Key，只有 Put / Get / Delete 接口，KV store 是线程安全的。

随机选择 10 个 key 做为事务的 key，Get 这 10 个 key 的 value, Sleep 100 毫秒，把其中 5 个 key 的 value 减 1，并把另外 5 个 key 的 value 加 1，Put 到 KV Store。

需要保证不能有两个事务同时修改一个 Key。

当所有事务执行停止的时候，要保证所有 value 之和是 0。

实现一个调度器，在并发数不限的条件下，让每分钟执行最多的事务。

需要完成可以运行的代码。


提示：
* 注意代码可读性，添加必要的注释（英文）
* 注意代码风格与规范，添加必要的单元测试和文档
* 注意异常处理，尝试优化性能


## 设计思路

KV store 为了保证线程安全，使用 sync.Map 代替。

sync.Map
version


## 测试

环境:
```
MacOS: 10.15.2 (19C57)
Processor Name:	Dual-Core Intel Core i7
Processor Speed: 1.7 GHz

go version go1.13.8 darwin/amd64
```

测试命令:
```bash
go run main.go
```

输出：
```bash

2020/02/23 11:42:28 Start test now.
2020/02/23 11:42:48 theads:   10, txn success QPS: 58.300000, total QPS: 97.000000, sum: 0, non zero count: 887
2020/02/23 11:43:08 theads:   50, txn success QPS: 135.500000, total QPS: 485.000000, sum: 0, non zero count: 928
2020/02/23 11:43:28 theads:  100, txn success QPS: 180.000000, total QPS: 975.000000, sum: 0, non zero count: 934
2020/02/23 11:43:48 theads:  200, txn success QPS: 222.350000, total QPS: 1950.000000, sum: 0, non zero count: 937
2020/02/23 11:44:08 theads:  500, txn success QPS: 263.650000, total QPS: 4937.600000, sum: 0, non zero count: 942
2020/02/23 11:44:28 theads: 1000, txn success QPS: 299.900000, total QPS: 9909.100000, sum: 0, non zero count: 947
2020/02/23 11:44:48 theads: 2000, txn success QPS: 336.000000, total QPS: 19717.700000, sum: 0, non zero count: 934
2020/02/23 11:45:08 theads: 3000, txn success QPS: 357.050000, total QPS: 29337.900000, sum: 0, non zero count: 941
2020/02/23 11:45:28 theads: 4000, txn success QPS: 370.250000, total QPS: 38840.300000, sum: 0, non zero count: 961
2020/02/23 11:45:49 theads: 5000, txn success QPS: 375.850000, total QPS: 47972.150000, sum: 0, non zero count: 947
2020/02/23 11:46:09 theads: 6000, txn success QPS: 371.350000, total QPS: 55009.200000, sum: 0, non zero count: 953
2020/02/23 11:46:29 theads: 7000, txn success QPS: 339.700000, total QPS: 56266.400000, sum: 0, non zero count: 947
2020/02/23 11:46:49 theads: 8000, txn success QPS: 304.050000, total QPS: 55155.600000, sum: 0, non zero count: 945
2020/02/23 11:47:09 theads: 9000, txn success QPS: 292.150000, total QPS: 55123.850000, sum: 0, non zero count: 943
2020/02/23 11:47:29 theads: 10000, txn success QPS: 279.950000, total QPS: 55048.450000, sum: 0, non zero count: 947

Process finished with exit code 0


```


sync.Map

```bash
sucess txn count: 5573, failed count: 701111
txn success QPS: 371.533333, sum: 0, non zero count: 941
```


## 题目描述

作业:

KV store 里有 1000 个 Key，只有 Put / Get / Delete 接口，KV store 是线程安全的。

随机选择 10 个 key 做为事务的 key，Get 这 10 个 key 的 value, Sleep 100 毫秒，把其中 5 个 key 的 value 减 1，并把另外 5 个 key 的 value 加 1，Put 到 KV Store。

需要保证不能有两个事务同时修改一个 Key。

当所有事务执行停止的时候，要保证所有 value 之和是 0。

实现一个调度器，在并发数不限的条件下，让每分钟执行最多的事务。

需要完成可以运行的代码。


提示：
* 注意代码可读性，添加必要的注释（英文）
* 注意代码风格与规范，添加必要的单元测试和文档
* 注意异常处理，尝试优化性能


