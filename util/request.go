package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func HTTPRequest[T any](url string, method string, params map[string]string) (response *T, err error) {
	kvset := []string{}
	for k, v := range params {
		kvset = append(kvset, k+"="+v)
	}
	param := strings.Join(kvset, "&")
	url = url + "?" + param

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
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
