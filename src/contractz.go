package main

func reqContract(z string) {
	println("request received... " + z)

	//senz := parse(z)

	// save event (request contract)
	t := eventTrans("restz", "orderzreq", "Contract request received")
	createTrans(t)

	// save event (broadcast contract)
	t = eventTrans("orderzreq", "*", "Broadcast contract request")
	createTrans(t)

	// publish to tranz
	kmsg := Kmsg{
		Topic: "tranz",
		Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
	}
	kchan <- kmsg

	// TODO find all chainz topics(designers, and printers) from etcd/zookeeper
	// TODO distribute contract to all chainz topics
	kmsg = Kmsg{
		Topic: "chainz",
		Msg:   z,
	}
	kchan <- kmsg
}

func respContract(z string) {
	println("response received... " + z)

	senz := parse(z)

	// save event (response received)
	t := eventTrans(senz.Sender, "orderzresp", "Contract response received")
	createTrans(t)

	// save event (send contract back)
	t = eventTrans("orderzreq", "restz", "Reply contract result")
	createTrans(t)

	// publish to tranz
	kmsg := Kmsg{
		Topic: "tranz",
		Msg:   tranzSenz(t.Id.String(), t.Type, t.Timestamp),
	}
	kchan <- kmsg

	// TODO handle contract responses from multiple chainz peers
	// TODO take decition according to the multiple chainz responses

	// TODO response back to restz topic
	kmsg = Kmsg{
		Topic: "restz",
		Msg:   z,
	}
	kchan <- kmsg
}
