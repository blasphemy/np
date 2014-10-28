package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Track struct {
	Album     string
	Artist    string
	Name      string
	PlayCount int
	tags      []string
}

func GetTrack(user string) Track {
	var tname string
	var aname string
	resp, err := http.Get(fmt.Sprintf("%s?method=user.getrecenttracks&user=%s&api_key=%s&format=json&limit=1", api_base_url, user, apikey))
	if err != nil {
		if debug {
			fmt.Println("Http Get 1:", err.Error())
		}
		os.Exit(1)
	}
	if debug {
		fmt.Println(resp)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			fmt.Println("Read resp.Body:", err.Error())
		}
		os.Exit(1)
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
		os.Exit(1)
	}
	if _, ok := data["recenttracks"].(map[string]interface{})["track"].([]interface{}); ok {
		aname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].([]interface{})[0].(map[string]interface{})["artist"].(map[string]interface{})["#text"].(string))
		tname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].([]interface{})[0].(map[string]interface{})["name"].(string))
	} else {
		aname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].(map[string]interface{})["artist"].(map[string]interface{})["#text"].(string))
		tname = url.QueryEscape(data["recenttracks"].(map[string]interface{})["track"].(map[string]interface{})["name"].(string))
	}

	resp, err3 := http.Get(fmt.Sprintf("%s?method=track.getInfo&api_key=%s&artist=%s&track=%s&format=json&username=%s&autocorrect=1", api_base_url, apikey, aname, tname, user))
	if err3 != nil {
		if debug {
			fmt.Println("Error in Get 2:", err3.Error())
		}
		os.Exit(1)
	}
	body, err3 = ioutil.ReadAll(resp.Body)
	if err3 != nil {
		if debug {
			fmt.Println("Error in Read 2:", err3.Error())
		}
		os.Exit(1)
	}
	var data4 map[string]interface{}
	err3 = json.Unmarshal(body, &data4)
	if err3 != nil {
		if debug {
			fmt.Println(err3.Error())
		}
		os.Exit(1)
	}
	rt := Track{}
	rt.Album = data4["track"].(map[string]interface{})["album"].(map[string]interface{})["title"].(string)
	rt.Artist = data4["track"].(map[string]interface{})["artist"].(map[string]interface{})["name"].(string)
	rt.Name = data4["track"].(map[string]interface{})["name"].(string)
	if _, ta := data4["track"].(map[string]interface{})["userplaycount"].(string); ta {
		rt.PlayCount, _ = strconv.Atoi(data4["track"].(map[string]interface{})["userplaycount"].(string))
	}
	if tags, lol := data4["track"].(map[string]interface{})["toptags"].(map[string]interface{}); lol {
		for _, j := range tags["tag"].([]interface{}) {
			rt.tags = append(rt.tags, j.(map[string]interface{})["name"].(string))
		}
	}
	return rt

}
