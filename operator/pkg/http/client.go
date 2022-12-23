package http

import (
	"encoding/json"
	"io"
	"net/http"
	"operators.kloudlite.io/pkg/errors"
)

func Get[T any](req *http.Request) (*T, *http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, resp, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp, errors.Newf("bad status code (%d), should be >= 200 & < 300", resp.StatusCode)
	}

	var b T
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	if err := json.Unmarshal(body, &b); err != nil {
		return nil, resp, err
	}
	return &b, resp, err
}
