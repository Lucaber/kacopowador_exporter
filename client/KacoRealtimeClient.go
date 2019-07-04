package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Realtime struct {
	Time time.Time
	PRealtime float64
}

type KacoRealtimeState struct {
	Realtime
}

func RequestRealtimeState(host string, date time.Time, ch chan<- KacoRealtimeState) {
	resp, err := http.Get(fmt.Sprintf("http://%s/realtime.csv?_=%d", host, date.Unix()))
	if err != nil {
		log.Println(err)
		ch <- KacoRealtimeState{Realtime{}}
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	realtime, err := parseRealtime(b)
	if err != nil {
		log.Println(err)
		ch <- KacoRealtimeState{Realtime{}}
		return
	}

	ch <- KacoRealtimeState{*realtime}
}

func parseRealtime(row []byte) (*Realtime, error) {
	values := bytes.Split(row, []byte(";"))

	if len(values) < 12 {
		return nil, &ParseError{"realtime row not found"}
	}

	ts, err := strconv.ParseInt(string(values[0]), 10, 64)
	if err != nil {
		panic(err)
	}
	Time := time.Unix(ts, 0)
	PRealtime, err := strconv.ParseFloat(string(values[len(values) - 3]), 64)
	PRealtime = PRealtime / (65535.0 / 100000.0) // from main.js
	if err != nil {
		return nil, err
	}
	return &Realtime{Time, PRealtime}, nil
}
