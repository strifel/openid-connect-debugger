package startconnection

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func Fetch(url string) (EndpointInformation, error) {
	resp, err := http.Get(url)

	if err != nil {
		return EndpointInformation{}, err
	}

	if resp.StatusCode != 200 {
		return EndpointInformation{}, errors.New("status code not 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EndpointInformation{}, err
	}

	var information EndpointInformation
	json.Unmarshal(body, &information)

	if information.Issuer == "" {
		return EndpointInformation{}, errors.New("issuer is empty")
	}

	if information.AuthorizationEndpoint == "" {
		return EndpointInformation{}, errors.New("authorization endpoint is empty")
	}

	if information.TokenEndpoint == "" {
		return EndpointInformation{}, errors.New("token endpoint is empty")
	}

	if information.UserinfoEndpoint == "" {
		return EndpointInformation{}, errors.New("userinfo endpoint is empty")
	}

	return information, nil
}
