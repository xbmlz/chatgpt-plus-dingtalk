package replicate

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/xbmlz/chatgpt-dingtalk/pkg/fetch"
)

type Replicate struct {
	BaseUrl  string
	ApiToken string
}

type ImageGenerateRequestInput struct {
	Prompt string `json:"prompt"`
}

type ImageGenerateRequest struct {
	Version string                    `json:"version"`
	Input   ImageGenerateRequestInput `json:"input"`
}

type ImageGenerateResponse struct {
	ID string `json:"id"`
}

type ImageGetResponse struct {
	Status string   `json:"status"`
	Output []string `json:"output"`
	Error  string   `json:"error"`
	// Urls   ImageGetResponseUrls   `json:"urls"`
}

func New(options Replicate) *Replicate {
	return &Replicate{
		BaseUrl:  options.BaseUrl,
		ApiToken: options.ApiToken,
	}
}

func (r *Replicate) Generate(param ImageGenerateRequest) (res string, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return res, err
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Token " + r.ApiToken,
	}
	generateBody, err := fetch.POST(r.BaseUrl+"/v1/predictions", headers, data)
	if err != nil {
		return res, err
	}

	generate := ImageGenerateResponse{}
	err = json.Unmarshal([]byte(generateBody), &generate)
	if err != nil {
		return res, err
	}
	get := ImageGetResponse{}
	// when status is succeeded && status !== 'failed'

	for {
		// sleep 1000s
		time.Sleep(1 * time.Second)
		getBody, err := fetch.GET(r.BaseUrl+"/v1/predictions/"+generate.ID, headers)
		if err != nil {
			return res, err
		}
		err = json.Unmarshal([]byte(getBody), &get)
		if err != nil {
			return res, err
		}
		if get.Status == "failed" {
			return res, errors.New(get.Error)
		}
		if get.Status == "succeeded" {
			// output[0].data
			res = get.Output[0]
			break
		}
	}
	return res, nil
}
