package tikaclient

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseUrl string
}

func NewClient(baseUrl string) Client {
	return Client{baseUrl: baseUrl}
}

func (client *Client) DetectLanguage(text string) (string, error) {
	requestUrl, err := url.JoinPath(client.baseUrl, "language/string")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPut, requestUrl, strings.NewReader(text))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "*/*")
	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(resBody), nil
}
