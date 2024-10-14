package decode

type DecodeRequest struct {
	InputString     string    `json:"inputString"`
}

type DecodeResponse struct {
	OutputString     string    `json:"outputString"`
}