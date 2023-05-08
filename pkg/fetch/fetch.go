package fetch

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

func GET(url string, headers map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", errors.New(string(body))
	}
	return string(body), nil
}

func POST(url string, headers map[string]string, data []byte) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", errors.New(string(body))
	}
	return string(body), nil
}
