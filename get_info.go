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
	data3 := &Loader{}
	err3 = json.Unmarshal(body, data3)
	if err3 != nil {
		if debug {
			fmt.Println(err3.Error())
		}
		os.Exit(1)
	}
	rt := RealTrack{}
	rt.Album = data3.Track_v.Album_v.Name_v
	rt.Artist = data3.Track_v.Artist_v.Name_2
	rt.Name = data3.Track_v.Name_v
	rt.PlayCount, _ = strconv.Atoi(data3.Track_v.UserCount)
	for _, k := range data3.Track_v.TopT.Tags {
		rt.tags = append(rt.tags, k.Name)
	}
	return rt
}
