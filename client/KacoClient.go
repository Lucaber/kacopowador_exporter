package client

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Info struct {
	KWh float64
}

type Current struct {
	Time time.Time
	Udc1 float64
	Idc1 float64
	Pdc1 float64
	Udc2 float64
	Idc2 float64
	Pdc2 float64
	Uac1 float64
	Iac1 float64
	Uac2 float64
	Iac2 float64
	Uac3 float64
	Iac3 float64
	Pdc  float64
	Pac  float64
	Tsys float64
}

type KacoState struct {
	Info
	Current
}


func RequestState(host string, date time.Time, ch chan<- KacoState) {
	resp, err := http.Get(fmt.Sprintf("http://%s/%d%02d%02d.CSV", host, date.Year(), date.Month(), date.Day()))
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		ch <- KacoState{Info{}, Current{}}
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	rows := bytes.Split(b, []byte("\r"))

	//File actually contains two csv "files"

	if len(rows) < 4 {
		ch <- KacoState{Info{}, Current{}}
		return
	}

	info, err := parseInfo(rows[:2])
	if err != nil {
		log.Fatal(err)
		ch <- KacoState{Info{}, Current{}}
		return
	}

	current, err := parseCurrent(rows[2:], date)
	if err != nil {
		log.Fatal(err)
		ch <- KacoState{Info{}, Current{}}
		return
	}

	ch <- KacoState{*info, *current}
}

func parseInfo(rows [][]byte) (*Info, error) {
	infoReader := csv.NewReader(bytes.NewReader(bytes.Join(rows, []byte("\n"))))
	infoReader.Comma = ';'

	infoRows, err := infoReader.ReadAll()
	if err != nil {
		return nil, err
	}
	KWh, err := strconv.ParseFloat(infoRows[1][4], 64)
	if err != nil {
		return nil, err
	}
	return &Info{KWh}, nil
}

func parseCurrent(rows [][]byte, date time.Time) (*Current, error) {
	infoReader := csv.NewReader(bytes.NewReader(bytes.Join(rows, []byte("\n"))))
	infoReader.Comma = ';'

	infoRows, err := infoReader.ReadAll()
	if err != nil {
		return nil, err
	}

	row := infoRows[len(infoRows)-1]
	Time, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", date.Format("2006-01-02"), row[0]))
	if err != nil {
		return nil, err
	}
	Udc1, err := strconv.ParseFloat(row[1], 64)
	Idc1, err := strconv.ParseFloat(row[2], 64)
	Pdc1, err := strconv.ParseFloat(row[3], 64)
	Udc2, err := strconv.ParseFloat(row[4], 64)
	Idc2, err := strconv.ParseFloat(row[5], 64)
	Pdc2, err := strconv.ParseFloat(row[6], 64)
	Uac1, err := strconv.ParseFloat(row[7], 64)
	Iac1, err := strconv.ParseFloat(row[8], 64)
	Uac2, err := strconv.ParseFloat(row[9], 64)
	Iac2, err := strconv.ParseFloat(row[10], 64)
	Uac3, err := strconv.ParseFloat(row[11], 64)
	Iac3, err := strconv.ParseFloat(row[12], 64)
	Pdc, err := strconv.ParseFloat(row[13], 64)
	Pac, err := strconv.ParseFloat(row[14], 64)
	Tsys, err := strconv.ParseFloat(row[15], 64)
	if err != nil {
		return nil, err
	}
	return &Current{Time, Udc1, Idc1, Pdc1, Udc2, Idc2, Pdc2, Uac1, Iac1, Uac2, Iac2, Uac3, Iac3, Pdc, Pac, Tsys}, nil
}
