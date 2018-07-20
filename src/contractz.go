package main

var rchans = make(map[string](chan string))

func reqContract(z string) {
	println("request received... " + z)

	// save event (request contract)
	t := eventTrans("rest", "opsreq", "Contract request received")
	createTrans(t)

	// save event (broadcast contract)
	t = eventTrans("opsreq", "*", "Broadcast contract request")
	createTrans(t)

	// publish to tranz
	kmsg := Kmsg{
		Topic: "tranz",
		Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
	}
	kchan <- kmsg

	// create channel and add to rchans with uuid
	c := make(chan string, 5)
	senz := parse(z)
	uid := senz.Attr["uid"]
	rchans[uid] = c

	// TODO find all chainz topics(designers, and printers) from etcd/zookeeper
	// TODO distribute contract to all chainz topics
	kmsg = Kmsg{
		Topic: "chain",
		Msg:   z,
	}
	kchan <- kmsg

	// wait for response
	waitForResponse(uid, c, 1)
}

func waitForResponse(uid string, c chan string, noPeers int) {
	var i int = 0
	responses := []string{}
	for {
		select {
		case r := <-c:
			println("reponse recived " + r)

			// append response
			responses = append(responses, r)

			i = i + 1
			if i == noPeers {
				// all peer responses received
				// TODO send response back to aws lambda
				println("all peers done....")

				// remove channel
				delete(rchans, uid)
			}
		}
	}
}

func respContract(z string) {
	println("response received... " + z)

	senz := parse(z)

	// save event (response received)
	t := eventTrans(senz.Sender, "opsresp", "Contract response received")
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
