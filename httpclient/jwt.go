package httpclient

import (
	"io/ioutil"
	"net/http"
	"st-go/errors"
	"encoding/json"
)

func GetJWTPublicKey(url string, v interface{}) error {
	client := http.Client{}
	req, err := http.NewRequest(`GET`, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.InternalServerError
	}

	return json.Unmarshal(body, &v)
}