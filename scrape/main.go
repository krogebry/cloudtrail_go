package main

import (
	"fmt"
	"path/filepath"
	"net/http"
	"net/url"
	"io/ioutil"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

func main() {

	//session, err := mgo.Dial("server1.example.com,server2.example.com")
	//if err != nil {
	//	panic(err)
	//}
	//defer session.Close()

	// dbconn := session.DB("cloudtrail")

	files, _ := filepath.Glob("/mnt/data/*/*")
	//fmt.Println(files)

	for _,filename := range files {
		fmt.Printf("File: %s\n", filename)
		resp, err := http.PostForm("http://192.168.56.100:8000/work",
			url.Values{"filename": {filename}})

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		fmt.Printf("Body: %s", body)
	}
}
