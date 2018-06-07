package main

import "fmt"
import "io/ioutil"
import "net/http"
import "bytes"
import "compress/gzip"

//import "encoding/json"
import "encoding/base64"

func DoPVerb(verb string, route string, params map[string]interface{}) string {

	//var buf, _ = json.Marshal(params)
	body := bytes.NewBuffer([]byte("name=test"))

	m := conf()
	token := m["token"]
	secret := m["secret"]
	prefix := m["url"]
	name := m["name"]
	url := fmt.Sprintf("%s/%s/%s", prefix, name, route)
	request, _ := http.NewRequest(verb, url, body)

	sEnc := base64.StdEncoding.EncodeToString([]byte(token + ":" + secret))

	request.Header.Set("Authorization", "BASIC "+sEnc)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/hal+json")
	client := &http.Client{}

	resp, err := client.Do(request)
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
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
	token := m["token"]
	secret := m["secret"]
	prefix := m["url"]
	name := m["name"]
	url := fmt.Sprintf("%s/%s/%s", prefix, name, route)
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
