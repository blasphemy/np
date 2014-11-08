package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

var (
	debug bool
)

func main() {
	viper.AddConfigPath("$HOME/np")
	viper.AddConfigPath("$HOME/.np")
	viper.AddConfigPath("$HOME")
	wd, err := os.Getwd()
	if err != nil {
		viper.AddConfigPath(wd)
	}
	viper.SetConfigFile("np")
	viper.SetDefault("template", "{{if .CurrentlyPlaying}}is now playing{{else}}last played{{end}} {{if .Artist}}[{{.Artist}}] - {{end}}{{if .Name}}[{{.Name}}] {{end}}{{if .Album}}On [{{.Album}}] {{end}}{{if .PlayCount}}[{{.PlayCount}} plays] {{end}}{{if .Tags}}[{{range $index, $element := .Tags}}{{if $index}} {{end}}#{{$element}}{{end}}]{{end}}")
	viper.SetDefault("debug", true)
	viper.SetDefault("username", "zxo0oxz")
	viper.ReadInConfig()
	fmt.Println(viper.GetString("test"))
	apikey := viper.GetString("apikey")
	username := viper.GetString("username")
	debug = viper.GetBool("debug")
	tmpl, err := template.New("format").Parse(viper.GetString("template"))
	if err != nil {
		panic(err)
	}
	track, err := GetTrack(username, apikey)
	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	err = tmpl.Execute(os.Stdout, track)
	if err != nil {
		panic(err)
	}
}
