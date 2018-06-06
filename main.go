package main

import "fmt"

import "github.com/codegangsta/cli"
import "time"
import "math/rand"
import "io/ioutil"

import "strings"
import "os"

//import "sort"

//import "github.com/antonholmquist/jason"

func main() {
	rand.Seed(time.Now().UnixNano())
	app := cli.NewApp()
	app.Name = "mql"
	app.Usage = "mode query language"
	app.Version = "16"
	app.Commands = []cli.Command{
		{Name: "reports", ShortName: "r",
			Usage: "reports", Action: ReportsAction},
		{Name: "token", ShortName: "t",
			Usage: "token", Action: TokenAction},
	}

	app.Run(os.Args)
}

func ReportsAction(c *cli.Context) {
	//email := c.Args().Get(0)
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
	fmt.Println(m["token"])

}
func TokenAction(c *cli.Context) {

}
