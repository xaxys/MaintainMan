package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func HTTPRequest[T any](url string, method string, params map[string]string, payload map[string]any) (response *T, err error) {
	kvset := []string{}
	for k, v := range params {
		kvset = append(kvset, k+"="+v)
	}
	param := strings.Join(kvset, "&")
	url = url + "?" + param

	client := &http.Client{}
	var bodyReader io.Reader
	if len(payload) > 0 {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	response = new(T)
	err = json.Unmarshal(body, response)
	return
}
