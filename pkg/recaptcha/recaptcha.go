package recaptcha

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ReCaptcha ..
type ReCaptcha struct {
	Secret string
}

// NewReCaptcha ..
func NewReCaptcha(secret string) *ReCaptcha {
	return &ReCaptcha{
		secret,
	}
}

type rcResp struct {
	Success bool
}

// Verify ..
func (rc *ReCaptcha) Verify(response, remoteip string) bool {
	params := url.Values{"secret": {rc.Secret}, "response": {response}, "remoteip": {remoteip}}
	resp, err := http.PostForm("https://www.recaptcha.net/recaptcha/api/siteverify", params)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var rp rcResp
	err = json.Unmarshal(body, &rp)
	if err != nil {
		return false
	}
	return rp.Success
}
