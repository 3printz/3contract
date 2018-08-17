package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type NotifyPr struct {
	AMC_ID      string
	PR_ID       string
	MatchStatus string
}

func notifyPreq(prId string, amcId string, status string) {
	// json req
	obj := NotifyPr{
		AMC_ID:      amcId,
		PR_ID:       prId,
		MatchStatus: status,
	}
	j, _ := json.Marshal(obj)

	log.Printf("INFO: send request to, %s with %s", apiConfig.prApi, string(j))

	// new request
	req, err := http.NewRequest("POST", apiConfig.prApi, bytes.NewBuffer(j))
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
