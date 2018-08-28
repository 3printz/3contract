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
		t := eventTrans("newco.biz", "newco.bcm", "Purchase request contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// handle purchase req, match amc/oem
		// create channel and add to rchans with uuid
		c := make(chan string, 5)
		uid := senz.Attr["uid"]
		prId := senz.Attr["prid"]
		rchans[uid] = c

		// TODO find all chainz topics(designers, and printers) from etcd/zookeeper
		topics := []string{"oem1", "amc1"}
		for _, topic := range topics {
			// save even
			t = eventTrans("newco.bcm", topic+".scm", "Purchase request contract")
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
		t := eventTrans("newco.biz", "newco.bcm", "Purchase order contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// handle purchase order
		// save even
		topics := []string{"oem1", "amc1"}
		for _, topic := range topics {
			// save even
			t = eventTrans("newco.bcm", topic+".scm", "Purchase order contract")
			createTrans(t)
		}

		notifyPorder(senz)
	}

	if senz.Attr["type"] == "DPREP" {
		// save event (request contract)
		t := eventTrans("newco.biz", "newco.bcm", "Data prep contract")
		createTrans(t)

		// publish to tranz
		kmsg := Kmsg{
			Topic: "tranz",
			Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
		}
		kchan <- kmsg

		// handle purchase order
		// save even
		topics := []string{"oem1", "amc1"}
		for _, topic := range topics {
			// save even
			t = eventTrans("newco.bcm", topic+".scm", "Data prep contract")
			createTrans(t)
		}

		notifyDprep(senz)
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
	t := eventTrans(senz.Sender+".scm", "newco.bcm", "Purchase request response")
	createTrans(t)

	// publish to tranz
	kmsg := Kmsg{
		Topic: "tranz",
		Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
	}
	kchan <- kmsg

	// find matching channel with uid and send z
	if c, ok := rchans[senz.Attr["uid"]]; ok {
		c <- z
	}
}
