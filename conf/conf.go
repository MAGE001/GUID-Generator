package conf

import (
	"flag"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	path string

	instOnce sync.Once
	inst     *config
)

func init() {
	flag.StringVar(&path, "conf", "", "conf path")
}

func New() *config {
	instOnce.Do(func() {
		flag.Parse()
		if path == "" {
			panic("invalid conf path")
		}
		inst = &config{}
		_, err := toml.DecodeFile(path, &inst)
		if err != nil {
			panic(err)
		}
	})
	return inst
}

type config struct {
	Listen    string
	NodeIdKey string
	Generator string
	Redis     struct {
		Addr             string
		Db               int
		PoolSize         int
		DialReadTimeout  int
		DialWriteTimeout int
		DialKeepAlive    int
		IdleTimeout      int
	}
}
