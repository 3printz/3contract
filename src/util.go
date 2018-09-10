package main

import (
	"strconv"
	"strings"
	"time"
)

func parse(msg string) Senz {
	fMsg := formatToParse(msg)
	tokens := strings.Split(fMsg, " ")
	senz := Senz{}
	senz.Msg = fMsg
	senz.Attr = map[string]string{}

	for i := 0; i < len(tokens); i++ {
		if i == 0 {
			senz.Ztype = tokens[i]
		} else if i == len(tokens)-1 {
			// signature at the end
			senz.Digsig = tokens[i]
		} else if strings.HasPrefix(tokens[i], "@") {
			// receiver @eranga
			senz.Receiver = tokens[i][1:]
		} else if strings.HasPrefix(tokens[i], "^") {
			// sender ^lakmal
			senz.Sender = tokens[i][1:]
		} else if strings.HasPrefix(tokens[i], "$") {
			// $key er2232
			key := tokens[i][1:]
			val := tokens[i+1]
			senz.Attr[key] = val
			i++
		} else if strings.HasPrefix(tokens[i], "#") {
			key := tokens[i][1:]
			nxt := tokens[i+1]

			if strings.HasPrefix(nxt, "#") || strings.HasPrefix(nxt, "$") ||
				strings.HasPrefix(nxt, "@") {
				// #lat #lon
				// #lat @eranga
				// #lat $key 32eewew
				senz.Attr[key] = ""
			} else {
				// #lat 3.2323 #lon 5.3434
				senz.Attr[key] = nxt
				i++
			}
		}
	}

	// set uid as the senz id
	senz.Uid = senz.Attr["uid"]

	return senz
}

func formatToParse(msg string) string {
	replacer := strings.NewReplacer(";", "", "\n", "", "\r", "")
	return strings.TrimSpace(replacer.Replace(msg))
}

func formatToSign(msg string) string {
	replacer := strings.NewReplacer(";", "", "\n", "", "\r", "", " ", "")
	return strings.TrimSpace(replacer.Replace(msg))
}

func uid() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	return config.senzieName + strconv.FormatInt(t, 10)
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func senzToTrans(senz *Senz) *Trans {
	trans := &Trans{}
	trans.Bank = config.senzieName
	trans.Id = uuid()
	trans.Timestamp = timestamp()
	trans.FromZaddress = senz.Sender
	trans.ToZaddress = senz.Receiver
	trans.PromizeBlob = "actual contract"
	trans.Type = senz.Attr["contract"]
	trans.Digsig = senz.Digsig

	return trans
}

func eventTrans(cid string, from string, to string, event string) *Trans {
	trans := &Trans{}
	trans.Bank = config.senzieName
	trans.Id = uuid()
	trans.Cid = cid
	trans.Timestamp = timestamp()
	trans.FromZaddress = from
	trans.ToZaddress = to
	trans.PromizeBlob = "actual contract"
	trans.Type = event

	return trans
}

func senzToUser(senz *Senz) *User {
	user := &User{}
	user.Zaddress = senz.Sender
	user.Bank = config.senzieName
	user.PublicKey = senz.Attr["pubkey"]
	user.Verified = false
	user.Active = true
	user.Timestamp = timestamp()

	return user
}

func regSenz() string {
	z := "SHARE #pubkey " + getIdRsaPubStr() +
		" #uid " + uid() +
		" @" + config.switchName +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func awaSenz(uid string) string {
	z := "AWA #uid " + uid +
		" @" + config.switchName +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func statusSenz(status string, uid string, to string) string {
	z := "DATA #status " + status +
		" #uid " + uid +
		" @" + to +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}

func tranzSenz(uid string, event string, time int64) string {
	z := "DATA #uid " + uid +
		" #time " + strconv.FormatInt(time, 10) +
		" #event " + event +
		" @" + "tranz" +
		" ^" + config.senzieName
	s, _ := sign(z, getIdRsa())

	return z + " " + s
}
