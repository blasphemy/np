package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetSpotifyUrl(i Track) string {
	url := BuildSpotifyQuery(i)
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return ""
	}
	if len(data["tracks"].(map[string]interface{})["items"].([]interface{})) < 1 {
		return ""
	}
	if r, ok := data["tracks"].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["external_urls"].(map[string]interface{})["spotify"].(string); ok {
		return r
	} else {
		return ""
	}
}

func BuildSpotifyQuery(i Track) string {
	baseurl := "https://api.spotify.com/v1/search?q="
	var fields []string
	if i.Name != "" {
		fields = append(fields, fmt.Sprintf("track:%s", url.QueryEscape(i.Name)))
	}
	if i.Artist != "" {
		fields = append(fields, fmt.Sprintf("artist:%s", url.QueryEscape(i.Artist)))
	}
	if i.Album != "" {
		fields = append(fields, fmt.Sprintf("album:%s", url.QueryEscape(i.Album)))
	}
	q := baseurl + strings.Join(fields, "+") + "&type=track"
	return q
}
