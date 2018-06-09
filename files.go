package main

import "io/ioutil"
import "runtime"
import "os"
import "strings"
import "strconv"

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func ReadList(name string) []string {
	data, _ := ioutil.ReadFile(UserHomeDir() + "/.mql_" + name)
	return strings.Split(string(data), ",")
}
func ReadLast(name string) int {
	data, _ := ioutil.ReadFile(UserHomeDir() + "/.mql_" + name + ".last")
	i, _ := strconv.Atoi(string(data))
	return i
}
func SaveLast(name string, index string) {
	ioutil.WriteFile(UserHomeDir()+"/.mql_"+name+".last", []byte(index), 0644)
}
func SaveList(name string, list []string) {
	ioutil.WriteFile(UserHomeDir()+"/.mql_"+name, []byte(strings.Join(list, ",")), 0644)
}
func SaveSQL(sql, token string) {
	ioutil.WriteFile(UserHomeDir()+"/.mql_"+token+".sql", []byte(sql), 0644)
}
func ReadSQL(token string) string {
	data, _ := ioutil.ReadFile(UserHomeDir() + "/.mql_" + token + ".sql")
	return string(data)
}
