package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type NotifyPr struct {
	Oem    string
	Amc    string
	Msg    string
	Status string
}

func notifyPreq() {
	// json req
	obj := NotifyPr{
		Oem:    "oem1",
		Amc:    "amc1",
		Msg:    "match done",
		Status: "SUCCESS",
	}
	j, _ := json.Marshal(obj)

	// new request
	uri := "http://dev.localhost" + apiConfig.prApi
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(j))
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: fail request, %s", err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	// status
	if resp.StatusCode != 200 {
		log.Printf("ERROR: fail request, status: %s response: %s", resp.StatusCode, string(b))
		return
	}
}

func notifyPorder() {

}

func notifyPrint() {

}
