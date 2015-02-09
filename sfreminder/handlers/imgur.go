package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"apikeys"
)

const imgurURL = `https://api.imgur.com/3/image`

type imgurRes struct {
	Data struct {
		Link string
	}
	Success bool
}

func copyToImgur(ctx appengine.Context, inputURL string) (string, error) {
	form := url.Values{}
	form.Add("image", inputURL)
	form.Add("type", "URL")
	req, err := http.NewRequest("POST", imgurURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", apikeys.ImgurKey))
	client := urlfetch.Client(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	var r imgurRes
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	if !r.Success  {
		ctx.Errorf("%V", r)
	}
	return r.Data.Link, nil
}
