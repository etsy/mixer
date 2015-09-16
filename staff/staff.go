package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	. "github.com/etsy/mixer/config"
	. "github.com/etsy/mixer/db"
)

func main() {
	err := Config.Load()
	if err != nil {
		log.Fatal("error reading or parsing config:", err)
	}

	log.Printf("populating staff tables from data feed\n")

	urls := Config.Staff.DatafeedUrl
	for _, url := range urls {
		/*log.Printf("url: %s\n", url)*/
		populateStaff(url)
	}

}

func populateStaff(url string) {

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	/*log.Printf("body is %s", body)*/

	var s []Staff
	if err := json.Unmarshal(body, &s); err != nil {
		log.Printf("ERROR: %v\n", err)
	}

	for _, v := range s {
		/*log.Printf("name: %s\n", v.FirstName)*/
		/*log.Printf("%#v\n", v)*/
		if v.Auth_UserName == "new_user_temp" || v.Auth_UserName == "" {
			continue
		}
		InsertStaffData(v)
	}
}
