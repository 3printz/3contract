package main

func reqContract(z string) {
	println("executing... " + z)

	// TODO save event (request received)

	// TODO save event (broadcast contract)

	// TODO find all chainz topics(designers, and printers) from etcd/zookeeper
	// TODO distribute contract to all chainz topics
	kmsg := Kmsg{
		Topic: "chainz",
		Msg:   z,
	}
	kchan <- kmsg
}

func respContract(z string) {
	println("received... " + z)

	// TODO save event

	// TODO handle contract responses from multiple chainz peers
	// TODO response back to restz topic
	kmsg := Kmsg{
		Topic: "restz",
		Msg:   z,
	}
	kchan <- kmsg
}
