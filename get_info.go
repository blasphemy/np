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
	Album     string
	Artist    string
	Name      string
	PlayCount int
	tags      []string
}

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
	if _, ok := data["recenttracks"].(map[string]interface{})["track"].([]interface{}); ok {
		aname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].([]interface{})[0].(map[string]interface{})["artist"].(map[string]interface{})["#text"].(string))
		tname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].([]interface{})[0].(map[string]interface{})["name"].(string))
	} else {
		aname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].(map[string]interface{})["artist"].(map[string]interface{})["#text"].(string))
		tname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].(map[string]interface{})["name"].(string))
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
				rt.tags = append(rt.tags, j.(map[string]interface{})["name"].(string))
			}
		} else {
			rt.tags = append(rt.tags, tags["tag"].(map[string]interface{})["name"].(string))
		}
	}
	return rt, nil
}
