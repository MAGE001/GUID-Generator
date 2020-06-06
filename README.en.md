[Language: [Original](./README.en.md)]

DISCLAIMER: Translated from Original using Google Translate.

# GUID-Generator
A distributed global unique id generator with two implementation: `twitter snowflake` and `/dev/urandom`.

## (1) twitter snowflake

#### ID Structure

```
/*
 * +------+----------------------+----------------+-----------+
 * | sign |     delta seconds    |    node id     | sequence  |
 * +------+----------------------+----------------+-----------+
 *   1bit          28bits              20bits         15bits
 */
```


#### Advantages

- High Performance
- Distributed
- Solves the problem of system clock callback

#### Disadvantages

- The entire distributed node supports restarting about 100W times（2^20）
- Time may lag in ID generation（When the request qps is low）Or ahead
（qps exceeded 2^15/s time）current system time，therefore, it is not suitable for scenarios that want to reflect business time by ID

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

snowflake generator uses redis storage worker node id，If you want to change the storage such as zookeeper, mysql, achive Storager ready interface。

1. Install Redis
2. Modify configuration file

```
NodeIdKey = "GUID:SNOWFLAKE:CURRENT-NODE-ID" # storage worker node id of key

Generator = "snowflake" # snowflake or random

[Redis]
Addr                = "127.0.0.1:6379"  # redis address
```
3. Start the service
```
$ ./httpsrv -conf ./httpserv.toml
```
4. Test
```
$ curl http://127.0.0.1:18080/ids?n=3
{"ids":[134423817710829569,134423817710829570,134423817710829571]}
```

## (2) Linux /dev/urandom

/dev/urandom - It can be used to generate cryptographically secure random numbers.

#### Advantages

- Distributed
- No external storage service required

#### Disadvantages

- Due to the existence of files I/O，Performance and snowflake difference 一orders of magnitude

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

#### Use

1. Modify configuration file
```
Generator = "random" # snowflake or random
```
2. Start service
```
$ ./httpsrv -conf ./httpserv.toml
```
3. Test
```
$ curl http://127.0.0.1:18080/ids?n=3
{"ids":[-7251799059149900519,8370896501513088158,7343472816840825557]}
```
