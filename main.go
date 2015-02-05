// hyapi project main.go
package main

import (
	"encoding/json"
	_ "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	lib *SicLib
)

func main() {
	port := flag.Int("port", 81, "http listen and serve the port,default=81")
	flag.Parse()

	lib = NewSicLib(`libSmartIndustryCode.dll`)
	defer lib.Free()
	lib.Init(1)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Write([]byte("error,only get"))
			return
		}
		desc := r.URL.Query().Get("d")
		result := Result{Codes: []ResultCode{}}
		codes := strings.Split(lib.GetCode(desc, 0, 5), "|")
		for _, oneCode := range strings.Split(codes[1], ";") {
			result.Codes = append(result.Codes, ResultCode{
				Code:     oneCode,
				Name:     HyCode[oneCode].Name,
				Category: HyCategory[HyCode[oneCode].Category],
			})
		}
		result.Hint = strings.Split(codes[2], ";")
		bys, err := json.Marshal(result)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`%s(%s);`, r.URL.Query().Get("callback"), bys)))
		return
	})
	fmt.Printf("start http://localhost:%d...\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
