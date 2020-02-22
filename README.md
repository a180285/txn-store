# txn-store

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


