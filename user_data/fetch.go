package userdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Fetch(url string, token string) (UserInfo, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UserInfo{}, "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)

	if err != nil {
		return UserInfo{}, "", err
	}

	if resp.StatusCode != 200 {
		return UserInfo{}, "", errors.New("HTTP " + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, "", err
	}

	var userInfo UserInfo
	json.Unmarshal(body, &userInfo)

	return userInfo, string(body), nil
}
