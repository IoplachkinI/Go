package server

import (
    "encoding/base64"
    "encoding/json"
    "math/rand/v2"
	"net/http"
	"time"
	"task2/decode"
)

func Decode(w http.ResponseWriter, req *http.Request) {
    reqBody := decode.DecodeRequest{}

    err := json.NewDecoder(req.Body).Decode(&reqBody)

    if err != nil {
        http.Error(w, "Couldn't parse request: " + err.Error(),
        http.StatusBadRequest)
        return
    }

    sDec, err := base64.StdEncoding.DecodeString(reqBody.InputString)
    if err != nil {
        http.Error(w, "Couldn't decode request string: " + err.Error(),
        http.StatusBadRequest)
        return
    }

    resp := decode.DecodeResponse{}
    resp.OutputString = string(sDec)
    
    body, err := json.Marshal(resp)

    if err != nil {
        http.Error(w, "Couldn't encode response string:" + err.Error(),
         http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(body)
}

func HardOp(w http.ResponseWriter, req *http.Request) {
    time.Sleep(time.Duration(rand.IntN(11) + 10) * time.Second)
    
    if rand.IntN(2) == 0 {
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        w.WriteHeader(http.StatusOK)
    }
}

func Version(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("v1.0.1"))
}
