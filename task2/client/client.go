package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"task2/decode"
)

func SendDecode(encMsg string, client *http.Client, url string) {
	reqBody := decode.DecodeRequest{encMsg}

	body, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal("Couldn't serialize request body:", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Couldn't create request:", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error parsing request:", err)
	}

	defer res.Body.Close()

	resp := &decode.DecodeResponse{}
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		log.Fatal("Couldn't parse response:", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal(res.Status)
	}

	fmt.Println("OutputString:", resp.OutputString)
}

func SendHardOp(client *http.Client, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Couldn't create request:" + err.Error())
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
  	defer cancel()
  	req = req.WithContext(ctx)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error parsing request:" + err.Error())
	}

	defer res.Body.Close()

	fmt.Println("Response code:", res.StatusCode)
}

func SendVersion(client *http.Client, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Couldn't create request:" + err.Error())
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error parsing request:" + err.Error())
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error reading request body:" + err.Error())
	}

	fmt.Println("Semver:", string(body))
}
