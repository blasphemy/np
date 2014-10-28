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

func GetTrack(user string) RealTrack {
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
	data := &Loader{}
	data2 := &Loader2{}
	err = json.Unmarshal(body, data)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		err2 := json.Unmarshal(body, data2)
		if err2 != nil {
			if debug {
				fmt.Println(err.Error())
			}
			os.Exit(1)
		}
	}
	if len(data2.Recent.TRay) > 0 {
		tname = url.QueryEscape(data2.Recent.TRay[0].Name_v)
		aname = url.QueryEscape(data2.Recent.TRay[0].Artist_v.Name_v)
	} else {
		tname = url.QueryEscape(data.Recent.Track_v.Name_v)
		aname = url.QueryEscape(data.Recent.Track_v.Artist_v.Name_v)
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
	//data3 := &Loader{}
	var data4 map[string]interface{}
	err3 = json.Unmarshal(body, &data4)
	if err3 != nil {
		if debug {
			fmt.Println(err3.Error())
		}
		os.Exit(1)
	}
	rt := RealTrack{}
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
	/*
		rt.Album = data3.Track_v.Album_v.Name_v
		for _, k := range data3.Track_v.TopT.Tags {
			rt.tags = append(rt.tags, k.Name)
		}
	*/
	return rt

}
