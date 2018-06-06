package main

import "fmt"

import "github.com/codegangsta/cli"
import "time"
import "math/rand"
import "io/ioutil"

import "strings"
import "os"

//import "sort"

import "github.com/antonholmquist/jason"

func main() {
	rand.Seed(time.Now().UnixNano())
	app := cli.NewApp()
	app.Name = "mql"
	app.Usage = "mode query language"
	app.Version = "16"
	app.Commands = []cli.Command{
		{Name: "spaces", ShortName: "s",
			Usage: "spaces", Action: SpacesAction},
		{Name: "reports", ShortName: "r",
			Usage: "reports", Action: ReportsAction},
		{Name: "token", ShortName: "t",
			Usage: "token", Action: TokenAction},
	}

	app.Run(os.Args)
}

func handleThing(thing, meta string) {
	v, _ := jason.NewObjectFromBytes([]byte(thing))
	if v == nil {
		return
	}
	e, _ := v.GetObject("_embedded")
	s, _ := e.GetObjectArray(meta)
	//token name
	for _, item := range s {
		stoken, _ := item.GetString("token")
		sname, _ := item.GetString("name")
		fmt.Println(stoken, sname)
	}
}
func conf() map[string]string {
	m := map[string]string{}
	b, _ := ioutil.ReadFile("conf/settings")
	prev := ""
	for i, line := range strings.Split(string(b), "\n") {
		if i%2 == 0 {
			prev = line
		} else {
			m[line] = prev
		}
	}
	return m
}
func SpacesAction(c *cli.Context) {
	//email := c.Args().Get(0)

	spaces := DoVerb("spaces")
	handleThing(spaces, "spaces")
}
func ReportsAction(c *cli.Context) {
	space_id := c.Args().Get(0)

	reports := DoVerb("spaces/" + space_id + "/reports")
	handleThing(reports, "reports")
}
func TokenAction(c *cli.Context) {

}
