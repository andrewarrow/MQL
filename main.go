package main

import "fmt"
import "strconv"
import "github.com/codegangsta/cli"
import "time"
import "math/rand"
import "io/ioutil"

import "strings"
import "os"
import "os/exec"

//import "sort"

import "github.com/antonholmquist/jason"

func main() {
	rand.Seed(time.Now().UnixNano())
	app := cli.NewApp()
	app.Name = "mql"
	app.Usage = "mode query language"
	app.Version = "16"
	app.Commands = []cli.Command{
		{Name: "spaces", ShortName: "p",
			Usage: "spaces", Action: SpacesAction},
		{Name: "reports", ShortName: "r",
			Usage: "reports", Action: ReportsAction},
		{Name: "queries", ShortName: "q",
			Usage: "queries", Action: QueriesAction},
		{Name: "sql", ShortName: "s",
			Usage: "sql", Action: SqlAction},
		{Name: "run", ShortName: "u",
			Usage: "run", Action: RunAction},
	}

	app.Run(os.Args)
}

func handleThing(thing, meta string, print bool) []*jason.Object {
	v, _ := jason.NewObjectFromBytes([]byte(thing))
	if v == nil {
		return []*jason.Object{}
	}
	e, _ := v.GetObject("_embedded")
	s, _ := e.GetObjectArray(meta)
	//token name
	list := []string{}
	for i, item := range s {
		stoken, _ := item.GetString("token")
		sname, _ := item.GetString("name")
		list = append(list, stoken)
		if print {
			fmt.Printf("%d. %s %s\n", i+1, stoken, sname)
		}
	}
	SaveList(meta, list)
	return s
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
	handleThing(spaces, "spaces", true)
}
func ReportsAction(c *cli.Context) {
	istr := c.Args().Get(0)
	i := ReadLast("space")
	if istr != "" {
		SaveLast("space", istr)
		i, _ = strconv.Atoi(istr)
	}
	list := ReadList("spaces")

	reports := DoVerb("spaces/" + list[i-1] + "/reports")
	handleThing(reports, "reports", true)
}
func QueriesAction(c *cli.Context) {
	istr := c.Args().Get(0)
	i := ReadLast("report")
	if istr != "" {
		SaveLast("report", istr)
		i, _ = strconv.Atoi(istr)
	}
	list := ReadList("reports")

	queries := DoVerb("reports/" + list[i-1] + "/queries")
	handleThing(queries, "queries", true)
}
func SqlAction(c *cli.Context) {
	i := ReadLast("report")
	j := ReadLast("query")
	jstr := c.Args().Get(0)
	if jstr != "" {
		SaveLast("query", jstr)
		j, _ = strconv.Atoi(jstr)
	}
	rlist := ReadList("reports")
	qlist := ReadList("queries")

	queries := DoVerb("reports/" + rlist[i-1] + "/queries")
	items := handleThing(queries, "queries", false)
	for _, item := range items {
		token, _ := item.GetString("token")
		sql, _ := item.GetString("raw_query")
		if token == qlist[j-1] {
			SaveSQL(sql)
		}
	}
	path := UserHomeDir() + "/.mql.sql"
	cmd := exec.Command("vim", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func handleLinks(thing, meta string, print bool) []*jason.Object {
	v, _ := jason.NewObjectFromBytes([]byte(thing))
	if v == nil {
		return []*jason.Object{}
	}
	e, _ := v.GetObject("_embedded")
	s, _ := e.GetObjectArray(meta)
	//token name
	for _, item := range s {
		l, _ := item.GetObject("_links")
		r, _ := l.GetObject("result")
		h, _ := r.GetString("href")
		if print {
			fmt.Println(h)
		}
	}
	return s
}
func RunAction(c *cli.Context) {
	report_id := c.Args().Get(0)
	query_id := c.Args().Get(1)
	//7fc4ef93285f/runs/a7b21118dbd8/query_runs/7a4122f98cbc/results
	r := DoVerb("reports/" + report_id + "/queries/" + query_id + "/runs")
	handleLinks(r, "query_runs", true)
}
