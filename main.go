package main

import "fmt"

//import "strconv"
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
	app.Name = "alamode"
	app.Usage = "alamode command line interface for mode"
	app.Version = "16"
	app.Commands = []cli.Command{
		{Name: "spaces", ShortName: "sp",
			Usage: "spaces", Action: SpacesAction},
		{Name: "reports", ShortName: "re",
			Usage: "reports", Action: ReportsAction},
		{Name: "queries", ShortName: "q",
			Usage: "queries", Action: QueriesAction},
		{Name: "sql", ShortName: "s",
			Usage: "sql", Action: SqlAction},
		{Name: "run", ShortName: "r",
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
	istr := c.Args().Get(0)
	if istr != "" {
		SaveLast("space", istr)
		return
	}
	spaces := DoVerb("spaces")
	handleThing(spaces, "spaces", true)
}
func ReportsAction(c *cli.Context) {
	istr := c.Args().Get(0)
	if istr != "" {
		SaveLast("report", istr)
		return
	}

	i := ReadLast("space")
	list := ReadList("spaces")

	/*if istr == "new" {
		params := map[string]interface{}{"name": "test"}
		DoPVerb("post", "spaces/"+list[i-1]+"/reports", params)
		return
	}*/

	reports := DoVerb("spaces/" + list[i-1] + "/reports")
	handleThing(reports, "reports", true)
}
func QueriesAction(c *cli.Context) {
	istr := c.Args().Get(0)
	if istr != "" {
		SaveLast("query", istr)
		return
	}

	i := ReadLast("report")
	list := ReadList("reports")

	queries := DoVerb("reports/" + list[i-1] + "/queries")
	handleThing(queries, "queries", true)
}
func RunAction(c *cli.Context) {
	i := ReadLast("report")
	j := ReadLast("query")
	rlist := ReadList("reports")
	qlist := ReadList("queries")

	sql := ReadSQL(qlist[j-1])
	thing := map[string]interface{}{"selected": false, "value": 100}
	rr := map[string]interface{}{"limit": thing}
	query := map[string]interface{}{"create_query_run": true,
		"limit": false, //"data_source_id": 8420,
		//"name": "People",
		"raw_query": sql, "token": qlist[j-1]}
	iqueries := []map[string]interface{}{query}

	report := map[string]interface{}{ //"name": "GunMeta", "description": "",
		"report_run": rr,
		"queries[]":  iqueries,
		"trk_source": "editor"}
	ireport := map[string]interface{}{"report": report}
	DoPVerb("POST", "reports/"+rlist[i-1]+"/runs", ireport)

	queries := DoVerb("reports/" + rlist[i-1] + "/queries")

	items := handleThing(queries, "queries", false)
	for _, item := range items {
		token, _ := item.GetString("token")
		dsi, _ := item.GetNumber("data_source_id")
		name, _ := item.GetString("name")
		if token == qlist[j-1] {
			r := DoVerb("reports/" + rlist[i-1] + "/queries/" + token + "/runs")
			handleLinks(r, "query_runs", false)
			SaveLast("query_run", "1")
			fmt.Println(dsi, name)
			break
		}
	}
	qlist = ReadList("query_runs")
	r := DoVerbFullPath(qlist[0])
	//fmt.Println(r)
	r = DoVerbFullPath(qlist[0] + "/content.json")
	fmt.Println(r)
}
func SqlAction(c *cli.Context) {
	i := ReadLast("report")
	j := ReadLast("query")
	rlist := ReadList("reports")
	qlist := ReadList("queries")

	queries := DoVerb("reports/" + rlist[i-1] + "/queries")
	items := handleThing(queries, "queries", false)
	for _, item := range items {
		token, _ := item.GetString("token")
		sql, _ := item.GetString("raw_query")
		if token == qlist[j-1] {
			SaveSQL(sql, token)
		}
	}
	path := UserHomeDir() + "/.mql_" + qlist[j-1] + ".sql"
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
	list := []string{}
	for i, item := range s {
		l, _ := item.GetObject("_links")
		r, _ := l.GetObject("result")
		href, _ := r.GetString("href")
		cra, _ := item.GetString("created_at")
		//coa, _ := item.GetString("completed_at")
		list = append(list, href)
		if print {
			tokens := strings.Split(cra, "T")
			fmt.Printf("%d. %s %s\n", i+1, tokens[0], strings.Split(tokens[1], ".")[0])
		}
	}
	SaveList(meta, list)
	return s
}
