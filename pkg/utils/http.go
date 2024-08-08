package utils

import "net/http"

func SendReq(r *http.Request) error {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
