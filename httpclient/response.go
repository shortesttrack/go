package httpclient

import (
	"net/http"
	"st-go/errors"
	"encoding/json"
	"io/ioutil"
)

type Response struct {
	*http.Response
}

func (r *Response) JSON(v interface{}) error {
	if r.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(r.StatusCode))
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	return err
}

func (r *Response) String() (string, error) {
	if r.StatusCode != http.StatusOK {
		return "", errors.New(http.StatusText(r.StatusCode))
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}