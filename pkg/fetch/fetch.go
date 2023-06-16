package fetch

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func GET(reqUrl string, headers map[string]string, timeout int, proxyUrl string) ([]byte, error) {
	return Request("GET", headers, reqUrl, nil, timeout, proxyUrl)
}

func POST(reqUrl string, headers map[string]string, data []byte, timeout int, proxyUrl string) ([]byte, error) {
	return Request("POST", headers, reqUrl, data, timeout, proxyUrl)
}

func DELETE(reqUrl string, headers map[string]string, timeout int, proxyUrl string) ([]byte, error) {
	return Request("DELETE", headers, reqUrl, nil, timeout, proxyUrl)
}

func Request(method string, headers map[string]string, reqUrl string, data []byte, timeout int, proxyUrl string) ([]byte, error) {
	req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	if proxyUrl != "" {
		uri, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, errors.New(string(body))
	}
	return body, nil
}
