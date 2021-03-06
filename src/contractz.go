package main

import (
	"log"
)

var rchans = make(map[string](chan string))

func reqContract(z string) {
	log.Printf("contract request received, %s", z)

	senz := parse(z)
	if senz.Attr["type"] == "PREQ" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "User accept Purchase Reqeust")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Create Purchase Request smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// handle purchase req, match amc/oem
		// create channel and add to rchans with uuid
		c := make(chan string, 5)
		uid := senz.Attr["uid"]
		prId := senz.Attr["prid"]
		rchans[uid] = c

		// TODO find all chainz topics(oems, amcs) from etcd/zookeeper
		topics := []string{"oem1", "amc1"}
		for _, topic := range topics {
			// save even
			t = eventTrans(senz.Attr["cid"], "newco.bcm", topic+".scm", "Broadcast PR smart contract")
			createTrans(t)

			// publish to kafka
			kmsg = Kmsg{
				Topic: topic,
				Msg:   z,
			}
			kchan <- kmsg
		}

		// wait for response
		waitForResponse(uid, prId, c, len(topics))
	}

	if senz.Attr["type"] == "PORD" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "Purchase order raised by User approved by NewCo")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Create Purchase Order smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// handle purchase order
		// save even
		topics := []string{"oem1", "amc1"}
		for _, topic := range topics {
			// save even
			t = eventTrans(senz.Attr["cid"], "newco.bcm", topic+".scm", "Notify Purchase Order contract")
			createTrans(t)
		}

		notifyPorder(senz)
	}

	if senz.Attr["type"] == "DPREP" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "Data preperation request")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Create Data Prep smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// handle purchase order
		// save even
		topics := []string{"oem1", "amc1"}
		for _, topic := range topics {
			// save even
			t = eventTrans(senz.Attr["cid"], "newco.bcm", topic+".scm", "Notify Data preperation contract")
			createTrans(t)
		}

		notifyDprep(senz)
	}

	if senz.Attr["type"] == "PRNT" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "Print request")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Create print smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// handle purchase order
		// save even
		topics := []string{"oem1", "amc1"}
		for _, topic := range topics {
			// save even
			t = eventTrans(senz.Attr["cid"], "newco.bcm", topic+".scm", "Notify Data preperation contract")
			createTrans(t)
		}
	}

	if senz.Attr["type"] == "DELNOTE" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "Delievey note request")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Create delivery note smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg
	}

	if senz.Attr["type"] == "INVOICE" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "Invoice request")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Create invoice smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// notify invoice
		notifyInvoice(senz)
	}

	if senz.Attr["type"] == "ACK" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "Ack request")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Ack smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg
	}

	if senz.Attr["type"] == "PAY" {
		// save event (request contract)
		t := eventTrans(senz.Attr["cid"], "newco.biz", "newco.bcm", "Payment request")
		createTrans(t)

		// save envet instantiate SC
		t = eventTrans(senz.Attr["cid"], "newco.bcm", "newco.scm", "Create payment smart contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// notify invoice
		notifyPayment(senz)
	}
}

func waitForResponse(uid string, prId string, c chan string, noPeers int) {
	var i int = 0
	responses := []string{}
	for {
		select {
		case r := <-c:
			log.Printf("response from peer: %s", r)

			// append response
			responses = append(responses, r)

			i = i + 1
			if i == noPeers {
				// all peer responses received, do matching logic
				// send response back
				log.Printf("all peers done uid: %s", "<UID>")
				for _, z := range responses {
					senz := parse(z)
					if senz.Attr["match"] == "YES" {
						notifyPreq(prId, senz.Attr["zid"], "SUCCESS")
					}
				}

				// remove channel at the end
				delete(rchans, uid)
			}
		}
	}
}

func respContract(z string) {
	log.Printf("contract response received, %s", z)

	senz := parse(z)

	// save event (response received)
	t := eventTrans(senz.Attr["cid"], senz.Sender+".scm", "newco.bcm", "Send PR response back")
	createTrans(t)

	// publish to tranz
	kmsg := Kmsg{
		Topic: "tranz",
		Msg:   tranzSenz(t.Cid, t.Type, t.Timestamp),
	}
	kchan <- kmsg

	// find matching channel with uid and send z
	if c, ok := rchans[senz.Attr["uid"]]; ok {
		c <- z
	}
}
