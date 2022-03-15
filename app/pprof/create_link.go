package pprof

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go_practicum/app/config"
	"io"
	"net/http"
	"strings"
)

func newRandomLink() string {
	return fmt.Sprintf("https://%s.com", uuid.NewString())
}

func CallCreateLinkHandler(client *http.Client, cfg config.Config) (string, error) {
	reqBody := strings.NewReader(newRandomLink())
	req, err := http.NewRequest(http.MethodPost, cfg.BaseURL, reqBody)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func CallCreateLinkJSONHandler(client *http.Client, cfg config.Config) (string, error) {
	type request struct {
		Link string `json:"url"`
	}

	type response struct {
		Output string `json:"result"`
	}

	reqBody, err := json.Marshal(request{
		Link: newRandomLink(),
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var respEncoded response
	err = json.NewDecoder(resp.Body).Decode(&respEncoded)
	if err != nil {
		return "", err
	}

	return respEncoded.Output, nil
}

func CallCreateAllLinkHandler(client *http.Client, cfg config.Config) (string, error) {
	type requestElem struct {
		CorrelationID string `json:"correlation_id"`
		Link          string `json:"original_url"`
	}

	type responseElem struct {
		CorrelationID string `json:"correlation_id"`
		Link          string `json:"short_url"`
	}

	inputs := &[]requestElem{
		{
			CorrelationID: "id1",
			Link:          newRandomLink(),
		},
		{
			CorrelationID: "id2",
			Link:          newRandomLink(),
		},
		{
			CorrelationID: "id3",
			Link:          newRandomLink(),
		},
	}
	reqBody, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.BaseURL+"/api/shorten/batch", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var respEncoded []responseElem
	err = json.NewDecoder(resp.Body).Decode(&respEncoded)
	if err != nil {
		return "", err
	}

	return respEncoded[0].Link, nil
}
