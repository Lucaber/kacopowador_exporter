package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Response struct {
	KWh string
}


func request(ch chan Response) {
	resp, err := http.Get("http://192.168.1.5/20190307.CSV")
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	b = bytes.ReplaceAll(b, []byte("\r"), []byte("\n"))
	reader := csv.NewReader(bytes.NewReader(b))
	reader.Comma = ';'

	lines, err := reader.ReadAll()
	if err != nil {
		println(err)
	}


	//str, err := ioutil.ReadAll(resp.Body)
	ch <- Response{KWh:strconv.Itoa(len(lines))}
}

func main() {
	ch := make(chan Response)

	c := cron.New()
	c.AddFunc("*/15 * * * * *", func() { request(ch) })
	c.Start()


	for {
		res := <-ch


		fmt.Println(res.KWh)
	}

}
