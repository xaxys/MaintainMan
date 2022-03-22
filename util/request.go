package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpRequest(url string, method string, params map[string]string, resjson any) error {
	kvset := []string{}
	for k, v := range params {
		kvset = append(kvset, k+"="+v)
	}
	param := strings.Join(kvset, "&")
	url = url + param

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &resjson)
	return nil
}
