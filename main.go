package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Loader2 struct {
	Recent RT2 `json:"recenttracks"`
}
type RT2 struct {
	TRay []Track `json:"track"`
}

type Loader struct {
	Recent  RT    `json:"recenttracks"`
	Track_v Track `json:"track"`
}

type RT struct {
	Track_v Track `json:"track"`
}

type Track struct {
	Name_v    string  `json:"name"`
	Artist_v  Artist  `json:"artist"`
	Album_v   Album   `json:"album"`
	UserCount string  `json:"userplaycount"`
	TopT      TopTags `json:"toptags"`
}

type Album struct {
	Name_v string `json:"title"`
}
type Artist struct {
	Name_v string `json:"#text"`
	Name_2 string `json:"name"`
}

type Tag struct {
	Name string `json:"name"`
}

type TopTags struct {
	Tags []Tag `json:"tag"`
}

var (
	username = "zxo0oxz"
	apikey string
	format   = "is listening to %s - %s On %s [%s plays]"
	debug    = true
)

const (
	api_base_url = "https://ws.audioscrobbler.com/2.0/"
)

func main() {
  apitemp, apierror := ioutil.ReadFile("api.config")
  if apierror != nil {
    if debug{
      fmt.Print("Error reading config file api.config")
      fmt.Print(apierror.Error())
      os.Exit(1)
    }
  }
  apikey = string(apitemp)
  apikey = strings.TrimSpace(apikey)
  resp, err := http.Get(fmt.Sprintf("%s?method=user.getrecenttracks&user=%s&api_key=%s&format=json&limit=1", api_base_url, username, apikey))
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	if debug {
		fmt.Println(resp)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	if debug {
		fmt.Println(string(body))
	}
	data := &Loader{}
	ndata := &Loader2{}
	err = json.Unmarshal(body, data)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		err = json.Unmarshal(body, ndata)
	}

	if debug {
		fmt.Println(data)
		fmt.Println(data.Recent.Track_v.Name_v)
		fmt.Println(data.Recent.Track_v.Artist_v.Name_v)
	}

	TRACK_NAME := url.QueryEscape(data.Recent.Track_v.Name_v)
	ARTIST_NAME := url.QueryEscape(data.Recent.Track_v.Artist_v.Name_v)

	if len(ndata.Recent.TRay) > 0 {

		TRACK_NAME = url.QueryEscape(ndata.Recent.TRay[0].Name_v)
		ARTIST_NAME = url.QueryEscape(ndata.Recent.TRay[0].Artist_v.Name_v)
	}

	resp, err = http.Get(fmt.Sprintf("%s?method=track.getInfo&api_key=%s&artist=%s&track=%s&format=json&username=%s&autocorrect=1", api_base_url, apikey, ARTIST_NAME, TRACK_NAME, username))
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	if debug {
		fmt.Println(string(body))
	}
	data2 := &Loader{}
	err = json.Unmarshal(body, data2)
	/*if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
  */
	if debug {
		fmt.Println(data2)
	}
	format = "is listening to [%s] - [%s]"
	o := fmt.Sprintf(format, data2.Track_v.Artist_v.Name_2, data2.Track_v.Name_v)
	if data2.Track_v.Album_v.Name_v != "" {
		o = o + fmt.Sprintf(" On [%s]", data2.Track_v.Album_v.Name_v)
	}
	if data2.Track_v.UserCount != "" {
		o = o + fmt.Sprintf(" [%s plays]", data2.Track_v.UserCount)
	}

	if len(data2.Track_v.TopT.Tags) > 0 {
		var k string
		for _, j := range data2.Track_v.TopT.Tags {
			k = k + fmt.Sprintf("#%s ", j.Name)
		}
		k = strings.TrimSpace(k)
		k = fmt.Sprintf(" [%s]", k)
		o = o + k
	}
	fmt.Print(o)
}
