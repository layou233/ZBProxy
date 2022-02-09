package version

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const Version = "3.0-rc.1"

func printErr(err error) {
	log.Printf("Error to check for update, caution: %v.", err.Error())
	log.Println(`You can check it yourself at https://github.com/layou233/ZBProxy/releases`)
}

func CheckUpdate() {
	resp, err := http.Get(`https://raw.githubusercontent.com/layou233/ZBProxy/master/version/version.go`)
	if err != nil {
		printErr(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printErr(err)
		return
	}
	if strings.Contains(string(body), Version) {
		fmt.Println("Your ZBProxy is up-to-date. Have fun!")
	} else {
		fmt.Println("Your ZBProxy is out of date! Check for the latest version at https://github.com/layou233/ZBProxy/releases")
	}
}
