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
	debug    = true
)

func main() {
	apitemp, apierror := ioutil.ReadFile("api.config")
	if apierror != nil {
		if debug {
			fmt.Println(apierror.Error())
		}
		os.Exit(1)
	}
	apikey = string(apitemp)
	apikey = strings.TrimSpace(apikey)

	track, err := GetTrack(username, apikey)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	format := "[%s] - [%s]"
	if track.CurrentlyPlaying {
		format = "is listening to " + format
	} else {
		format = "last played " + format
	}
	o := fmt.Sprintf(format, track.Artist, track.Name)
	if track.Album != "" {
		o = o + fmt.Sprintf(" On [%s]", track.Album)
	}
	if track.PlayCount > 0 {
		o = o + fmt.Sprintf(" [%d plays]", track.PlayCount)
	}
	if len(track.Tags) > 0 {
		var k string
		for _, j := range track.Tags {
			k = k + fmt.Sprintf("#%s ", strings.Replace(j, " ", "", -1))
		}
		k = strings.TrimSpace(k)
		k = fmt.Sprintf(" [%s]", k)
		o = o + k
	}
	fmt.Print(o)
}
