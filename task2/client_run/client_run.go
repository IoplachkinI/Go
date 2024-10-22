package main

import (
	"encoding/base64"
	"net/http"
	cl "task2/client"
)

func main() {
	client := &http.Client{}

	msg := "Awesome test message"
	encMsg := base64.StdEncoding.EncodeToString([]byte(msg))

	cl.SendDecode(encMsg, client, "http://localhost:8080/decode")

	cl.SendHardOp(client, "http://localhost:8080/hard-op")

	cl.SendVersion(client, "http://localhost:8080/version")
}