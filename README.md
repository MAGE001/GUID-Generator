# GUID-Generator
A distributed global unique id generator with two implementation: `twitter snowflake` and `/dev/urandom`.

## (1) twitter snowflake

#### ID结构

```
/*
 * +------+----------------------+----------------+-----------+
 * | sign |     delta seconds    |    node id     | sequence  |
 * +------+----------------------+----------------+-----------+
 *   1bit          28bits              20bits         15bits
 */
```


#### 优点

- 高性能
- 分布式
- 解决系统时钟回拨问题

#### 缺点

- 整个分布式节点支持重启约100W次（2^20）
- 生成ID中时间可能滞后（请求qps较低时）或超前（qps超过2^15/s时）于系统当前时间，因此不适用于希望通过ID反映业务时间的场景

#### benchmark

```
$ go test -bench=. -run=none
goos: darwin
goarch: amd64
pkg: github.com/GUID-Generator/snowflake
Benchmark_10000-8   	    2000	    957087 ns/op
Benchmark_1000-8    	   20000	     90390 ns/op
Benchmark_100-8     	  200000	      8728 ns/op
Benchmark_1-8       	20000000	       108 ns/op
PASS
ok  	github.com/GUID-Generator/snowflake	8.886s
```

#### 使用

snowflake generator使用redis存储worker node id，如果想换成其它存储如zookeeper、mysql，实现Storager即可接口。

1. 安装redis
2. 修改配置文件
```
NodeIdKey = "GUID:SNOWFLAKE:CURRENT-NODE-ID" # 存储worker node id的key

Generator = "snowflake" # snowflake or random

[Redis]
Addr                = "127.0.0.1:6379"  # redis地址
```
3. 启动服务
```
$ ./httpsrv -conf ./httpserv.toml
```
4. 测试
```
$ curl http://127.0.0.1:18080/ids?n=3
{"ids":[134423817710829569,134423817710829570,134423817710829571]}
```

## (2) Linux /dev/urandom

/dev/urandom可以用来生成密码学上的安全随机数。

#### 优点

- 分布式
- 不需要外部存储辅助

#### 缺点

- 由于存在文件I/O，性能与snowflake相差一个数量级

#### benchmark

```
$ go test -bench=. -run=none
goos: darwin
goarch: amd64
pkg: github.com/GUID-Generator/random
Benchmark_10000-8   	     200	   6413602 ns/op
Benchmark_1000-8    	    2000	    633439 ns/op
Benchmark_100-8     	   20000	     63238 ns/op
Benchmark_1-8       	 2000000	       648 ns/op
PASS
ok  	github.com/GUID-Generator/random	7.151s
```

#### 使用

1. 修改配置文件
```
Generator = "random" # snowflake or random
```
2. 启动服务
```
$ ./httpsrv -conf ./httpserv.toml
```
3. 测试
```
$ curl http://127.0.0.1:18080/ids?n=3
{"ids":[-7251799059149900519,8370896501513088158,7343472816840825557]}
```
