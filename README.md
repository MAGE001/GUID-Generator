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

- 整个分布式节点支持有限次重启 2^20 约100W次
- 生成ID中时间可能滞后（请求qps较低时）或超前（qps超过2^15/s时）于系统当前时间，因此不适用通过ID反应业务时间的场景。也因此会导致服务重启后，生成的ID与前一次之间有较大gap。

#### Benchmark

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


## (2) Linux /dev/urandom

#### 优点

- 不需要外部存储辅助

#### 缺点

- 由于存在文件I/O，性能与snowflake相差一个数量级

#### Benchmark

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
