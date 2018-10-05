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

type NotifyAmc struct {
	PO_ID  string
	AMC_ID string
}

type NotifyOem struct {
	PO_ID  string
	OEM_ID string
}

type NotifyDprep struct {
	PO_ID     string
	ENTITY_ID string
}

type NotifyInvoice struct {
	INVOICE_ID string
	ENTITY_ID  string
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

func notifyPorder(senz Senz) {
	// notify amc
	obj1 := NotifyAmc{
		PO_ID:  senz.Attr["poid"],
		AMC_ID: senz.Attr["amcid"],
	}
	j, _ := json.Marshal(obj1)
	notify(j, senz.Attr["amcapi"])

	// notify oem
	obj2 := NotifyOem{
		PO_ID:  senz.Attr["poid"],
		OEM_ID: senz.Attr["oemid"],
	}
	j, _ = json.Marshal(obj2)
	notify(j, senz.Attr["oemapi"])
}

func notifyDprep(senz Senz) {
	// notify amc
	obj1 := NotifyDprep{
		PO_ID:     senz.Attr["poid"],
		ENTITY_ID: senz.Attr["amcid"],
	}
	j, _ := json.Marshal(obj1)
	notify(j, senz.Attr["amcapi"])
}

func notifyInvoice(senz Senz) {
	// notify amc
	obj1 := NotifyInvoice{
		INVOICE_ID: senz.Attr["inid"],
		ENTITY_ID:  senz.Attr["customerId"],
	}
	j, _ := json.Marshal(obj1)
	notify(j, senz.Attr["callback"])
}

func notifyPayment(senz Senz) {
	// notify amc
	obj1 := NotifyInvoice{
		INVOICE_ID: senz.Attr["inid"],
		ENTITY_ID:  senz.Attr["entityId"],
	}
	j, _ := json.Marshal(obj1)
	notify(j, senz.Attr["callback"])
}

func notify(j []byte, uri string) {
	log.Printf("INFO: send request to, %s with %s", uri, string(j))

	// new request
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
