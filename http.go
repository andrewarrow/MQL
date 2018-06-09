package main

import "fmt"
import "io/ioutil"
import "net/http"

//import "net/url"
import "bytes"
import "compress/gzip"

import "encoding/json"
import "encoding/base64"

func DoPVerb(verb string, route string, params map[string]interface{}) string {

	var buf, _ = json.Marshal(params)
	fmt.Println(string(buf))
	//body := bytes.NewBuffer([]byte("raw_query=" + url.QueryEscape(sql)))
	body := bytes.NewBuffer(buf)

	m := conf()
	token := m["token"]
	secret := m["secret"]
	prefix := m["url"]
	name := m["name"]
	url := fmt.Sprintf("%s/api/%s/%s", prefix, name, route)
	request, _ := http.NewRequest(verb, url, body)

	sEnc := base64.StdEncoding.EncodeToString([]byte(token + ":" + secret))

	request.Header.Set("Authorization", "BASIC "+sEnc)
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cache-Control", "no-cache")
	// 	request.Header.Set("Accept", "application/hal+json")
	// 'Cache-Control': 'no-cache',
	client := &http.Client{}

	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			if resp.StatusCode == 202 {
				return string(body)
			} else {
				fmt.Println(string(body))
			}
		} else {
			fmt.Println(string(body), err)
		}
	} else {
		fmt.Println(err)
	}
	return ""
}

func DoVerbFullPath(route string) string {
	m := conf()
	token := m["token"]
	secret := m["secret"]
	prefix := m["url"]
	url := fmt.Sprintf("%s/%s", prefix, route)
	request, _ := http.NewRequest("GET", url, nil)

	sEnc := base64.StdEncoding.EncodeToString([]byte(token + ":" + secret))

	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Set("Authorization", "BASIC "+sEnc)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/hal+json")
	client := &http.Client{}

	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		reader, err := gzip.NewReader(resp.Body)
		body, err := ioutil.ReadAll(reader)
		if err == nil {
			if resp.StatusCode == 200 {
				return string(body)
			} else {
				fmt.Println(string(body))
			}
		} else {
			fmt.Println(string(body), err)
		}
	} else {
		fmt.Println(err)
	}
	return ""
}
func DoVerb(route string) string {
	m := conf()
	name := m["name"]
	return DoVerbFullPath("api/" + name + "/" + route)
}
