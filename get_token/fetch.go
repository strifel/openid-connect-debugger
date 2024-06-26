package gettoken

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	nurl "net/url"
)

func Fetch(url string, id string, secret string, code string, callbackEndpoint string) (CodeResponse, error) {
	resp, err := http.PostForm(url, nurl.Values{
		"client_id":     {id},
		"client_secret": {secret},
		"redirect_uri":  {"http://" + callbackEndpoint + ":8070/callback"},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	})

	if err != nil {
		return CodeResponse{}, err
	}

	if resp.StatusCode != 200 {
		return CodeResponse{}, errors.New("HTTP " + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CodeResponse{}, err
	}

	var response CodeResponse
	json.Unmarshal(body, &response)

	// We will only use AccessToken going forward,
	// so we will only check if it exists
	if response.AccessToken == "" {
		return response, errors.New("access token is empty")
	}

	return response, nil
}
