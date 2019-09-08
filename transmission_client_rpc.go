package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TrClient Transmission RPC Client
type TrClient struct {
	client *http.Client
}

// TrResponse Transmission RPC response
type TrResponse struct {
	Result    string          `json:"result"`
	Arguments json.RawMessage `json:"arguments"`
	Tag       int             `json:"tag"`
}

// TrRequest Transmission RPC request
type TrRequest struct {
	Method    string          `json:"method"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
	Tag       *int            `json:"tag,omitempty"`
}

type TransmissionError struct {
	Code   int
	Detail string
}

func (e *TransmissionError) Error() string {
	return fmt.Sprintf("Transmission Error: Code=%d, Detail=%s\n", e.Code, e.Detail)
}

// NewTrClient constructor for Transmission Client
func NewTrClient() TrClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	return TrClient{client}
}

func (t *TrClient) _sendRequest(req *http.Request) (*http.Response, error) {
	res, err := t.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	if res.StatusCode == 409 {
		return res, &TransmissionError{res.StatusCode, "Refresh Session ID"}
	}

	return res, nil

}

func (t *TrClient) _makeRequest(params TrRequest) (*http.Request, error) {
	url := fmt.Sprintf(TransmissionRPCUrl, TransmissionRPCServerName)
	buf := new(bytes.Buffer)
	fmt.Printf("Params -> %v\n", params)
	json.NewEncoder(buf).Encode(params)

	fmt.Printf("Buf -> %s\n", buf)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Println("New request generate failed")
		return &http.Request{}, err
	}
	fmt.Printf("Generated Request -> %s \n", req.Body)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transmission-Session-Id", TransmissionSessionID)
	fmt.Printf("Session %s\n", TransmissionSessionID)
	return req, nil
}

func (t *TrClient) makeRequest(params TrRequest) (*TrResponse, error) {
	req, err := t._makeRequest(params)

	if err != nil {
		return &TrResponse{}, err
	}

	res, err := t._sendRequest(req)

	if err != nil {
		switch err.(type) {
		case *TransmissionError:
			TransmissionSessionID = res.Header.Get("X-Transmission-Session-Id")
			fmt.Printf("SESSION ID -> %s\n", TransmissionSessionID)
			req, err = t._makeRequest(params)

			if err != nil {
				return &TrResponse{}, err
			}

			res, err = t._sendRequest(req)

			if err != nil {
				return &TrResponse{}, err
			}

		default:
			return &TrResponse{}, err
		}
	}

	defer res.Body.Close()

	r := &TrResponse{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(r)

	if err != nil {
		return &TrResponse{}, err
	}

	return r, nil
}

func (t *TrRequest) String() string {
	return fmt.Sprintf(
		"Method: %v\nArguments: %v\nTag: %v\n", t.Method, t.Arguments, *t.Tag,
	)
}
