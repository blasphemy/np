package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	username = "zxo0oxz"
	apikey   string
	format   = "is listening to %s - %s On %s [%s plays]"
	debug    = false
)

const (
	api_base_url = "https://ws.audioscrobbler.com/2.0/"
)

func main() {
	apitemp, apierror := ioutil.ReadFile("/home/daniel/np/api.config")
	if apierror != nil {
		if debug {
			fmt.Println(apierror.Error())
		}
		os.Exit(1)
	}
	apikey = string(apitemp)
	apikey = strings.TrimSpace(apikey)

	track := GetTrack(username)
	format = "is listening to [%s] - [%s]"
	o := fmt.Sprintf(format, track.Artist, track.Name)
	if track.Album != "" {
		o = o + fmt.Sprintf(" On [%s]", track.Album)
	}
	if track.PlayCount > 0 {
		o = o + fmt.Sprintf(" [%d plays]", track.PlayCount)
	}
	if len(track.tags) > 0 {
		var k string
		for _, j := range track.tags {
			k = k + fmt.Sprintf("#%s ", j)
		}
		k = strings.TrimSpace(k)
		k = fmt.Sprintf(" [%s]", k)
		o = o + k
	}
	fmt.Print(o)
}
