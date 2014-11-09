package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	api_base_url = "https://ws.audioscrobbler.com/2.0/"
)

type Track struct {
	Album            string
	Artist           string
	Name             string
	PlayCount        int
	Tags             []string
	CurrentlyPlaying bool
	SpotifyUrl       string
}

/* GetTrack returns a Track. One or all fields may be nil, because last.fm rocks
   This function will eventually be split off into it's own library.
*/
func GetTrack(user string, ApiKey string) (Track, error) {
	var tname string
	var aname string
	var resp *http.Response
	var err error
	rt := Track{}
	for i := 1; i < 3; i++ {
		if debug {
			fmt.Println("HTTP GET 1 TRY ", i)
		}
		resp, err = http.Get(fmt.Sprintf("%s?method=user.getrecenttracks&user=%s&api_key=%s&format=json&limit=1", api_base_url, user, ApiKey))
		if err == nil {
			break
		}
	}
	if err != nil {
		if debug {
			fmt.Println("Http Get 1:", err.Error())
		}
		return rt, err
	}
	if debug {
		fmt.Println(resp)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			fmt.Println("Read resp.Body:", err.Error())
		}
		return rt, err
	}
	if debug {
		fmt.Println(string(body))
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		return rt, err
	}
	var ct map[string]interface{}
	if _, ok := data["recenttracks"].(map[string]interface{})["track"].([]interface{}); ok {
		ct = data["recenttracks"].(map[string]interface{})["track"].([]interface{})[0].(map[string]interface{})
	} else {
		ct = data["recenttracks"].(map[string]interface{})["track"].(map[string]interface{})
	}
	aname = url.QueryEscape(ct["artist"].(map[string]interface{})["#text"].(string))
	tname = url.QueryEscape(ct["name"].(string))
	if cc, ok := ct["@attr"].(map[string]interface{}); ok {
		if aa, bb := cc["nowplaying"].(string); bb && aa == "true" {
			rt.CurrentlyPlaying = true
			if debug {
				fmt.Println("Now playing = true")
			}
		} else {
			rt.CurrentlyPlaying = false
			fmt.Println("Now playing = false")
		}
	}

	for i := 1; i < 3; i++ {
		if debug {
			fmt.Println("HTTP GET 2 TRY ", i)
		}
		resp, err = http.Get(fmt.Sprintf("%s?method=track.getInfo&api_key=%s&artist=%s&track=%s&format=json&username=%s&autocorrect=1", api_base_url, ApiKey, aname, tname, user))
		if err == nil {
			break
		}
	}
	if err != nil {
		if debug {
			fmt.Println("Error in Get 2:", err.Error())
		}
		return rt, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			fmt.Println("Error in Read 2:", err.Error())
		}
		return rt, err
	}
	if debug {
		fmt.Println(string(body))
	}
	var data2 map[string]interface{}
	err = json.Unmarshal(body, &data2)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		return rt, err
	}
	if al, kk := data2["track"].(map[string]interface{})["album"].(map[string]interface{}); kk {
		rt.Album = al["title"].(string)
	}
	rt.Artist = data2["track"].(map[string]interface{})["artist"].(map[string]interface{})["name"].(string)
	rt.Name = data2["track"].(map[string]interface{})["name"].(string)
	if _, ta := data2["track"].(map[string]interface{})["userplaycount"].(string); ta {
		rt.PlayCount, _ = strconv.Atoi(data2["track"].(map[string]interface{})["userplaycount"].(string))
	}
	if tags, lol := data2["track"].(map[string]interface{})["toptags"].(map[string]interface{}); lol {
		if t, lol2 := tags["tag"].([]interface{}); lol2 {
			for _, j := range t {
				rt.Tags = append(rt.Tags, j.(map[string]interface{})["name"].(string))
			}
		} else {
			rt.Tags = append(rt.Tags, tags["tag"].(map[string]interface{})["name"].(string))
		}
	}
	rt.SpotifyUrl = GetSpotifyUrl(rt)
	return rt, nil
}
