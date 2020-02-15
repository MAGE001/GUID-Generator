package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/GUID-Generator/conf"
	"github.com/GUID-Generator/generator"
	"github.com/GUID-Generator/random"
	"github.com/GUID-Generator/snowflake"
	"github.com/GUID-Generator/storage"
)

func main() {
	httpServe()
}

func httpServe() {
	var gen generator.Generator
	if conf.New().Generator == "snowflake" {
		gen = snowflake.NewSnowflakeGenerator(storage.NewRStorage())
	} else if conf.New().Generator == "random" {
		gen = random.NewRandomGenerator()
	} else {
		panic(fmt.Sprintf("invalid generator: %s", conf.New().Generator))
	}

	http.HandleFunc("/ids", func(w http.ResponseWriter, r *http.Request) {
		nstr := r.URL.Query().Get("n")
		if nstr == "" {
			http.Error(w, "lack of param", http.StatusBadRequest)
			return
		}

		n, err := strconv.Atoi(nstr)
		if err != nil || n <= 0 {
			http.Error(w, "invalid param", http.StatusBadRequest)
			return
		}

		resp := &response{
			Ids: gen.NextIds(n),
		}
		t, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(t))
	})

	log.Fatal(http.ListenAndServe(conf.New().Listen, nil))
}

type response struct {
	Ids []int64 `json:"ids"`
}
