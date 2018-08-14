package main

import (
	"log"
)

var rchans = make(map[string](chan string))

func reqContract(z string) {
	log.Printf("contract request received, %s", z)

	// save event (request contract)
	t := eventTrans("3rest", "3ops", "Contract request received")
	createTrans(t)

	// publish to tranz
	kmsg := Kmsg{
		Topic: "tranz",
		Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
	}
	kchan <- kmsg

	senz := parse(z)
	if senz.Attr["type"] == "PREQ" {
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
			t = eventTrans("3ops", topic, "Send contract request")
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
		// handle purchase order
		// save even
		topic := senz.Attr["oemid"]
		t = eventTrans("3ops", topic, "Send contract request")
		createTrans(t)

		// call oem to get design via kafka
		kmsg = Kmsg{
			Topic: topic,
			Msg:   z,
		}
		kchan <- kmsg
	}

	if senz.Attr["type"] == "PRINT" {
		// handle print request
		// call amc to print
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
				// all peer responses received, send response back
				log.Printf("all peers done uid: %s", "<UID>")
				notifyPreq(prId)

				// remove channel
				delete(rchans, uid)
			}
		}
	}
}

func respContract(z string) {
	log.Printf("contract response received, %s", z)

	senz := parse(z)

	// save event (response received)
	t := eventTrans(senz.Sender, "3ops", "Contract response received")
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
