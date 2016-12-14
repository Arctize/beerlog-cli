/* See LICENSE file for copyright and license details.
 *
 * This program is a commandline frontend for the AMIV Bierlog server.
 * It works by querying the said server and parsing the HTML.
 *
 * AUTHOR: Patrick Wicki <patrick.wicki@fsfe.org>
 */

package main

import (
	"flag"
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
	"strconv"
	"strings"
	"os"
)

func main(){

	cur_time := time.Now().Local().Format("2006-01-02")

	// Parse cli arguments
	date_start := flag.String("sdate", "", "Start date")
	date_end := flag.String("edate", cur_time , "End date")

	org := flag.String("org", "amiv" , "Department")
	userid := flag.String("user", "", "User ID, leave empty to show all")

	size := flag.Int("size", 10, "Number of queries to show")
	tp := flag.String("type", "coffee", "Type to query, can be one of: {beer,coffee,all}")

	flag.Parse()

	// Create the GET request with url encoding
	url := "http://intern.amiv.ethz.ch/beerlog/index.php"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Encode query url after adding parameters
	q := req.URL.Query()
	q.Add("date1", *date_start)
	q.Add("date2", *date_end)
	q.Add("org", *org)
	q.Add("userid", *userid)
	q.Add("size", strconv.Itoa(*size))
	q.Add("type", *tp)
	req.URL.RawQuery = q.Encode()

	// Define client and use it to query the beerlog server
	fmt.Printf("Querying server... ")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("response Status: %s\n\n", resp.Status)
		fmt.Println("Error connecing to the server.")
		resp.Body.Close()
	}
	defer os.Exit(1)

	fmt.Printf("response Status: %s\n\n", resp.Status)
	
	// Read body and convert it into a string
	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	if strings.Count(content, "td") == 0 {
		fmt.Println("No results.")
	}
	defer os.Exit(0)

	// TODO: This part sucks. Proper HTML parsing should be done here.
	tbody := content[strings.Index(content, "<td>"): strings.LastIndex(content,"</td>")+1]
	lines := strings.Split(tbody, " ")

	separator := strings.Repeat("-", 54)
	fmt.Println(separator)
	fmt.Printf("%6s%16s%16s%16s\n","#","User","Organisation","Consumption")
	fmt.Println(separator)

	for i := 0; i < len(lines); i++ {
		// TODO: This part sucks even more.
		line := strings.Split(lines[i], "</td><td>")
		rank := line[0][strings.Index(line[0], ">")+1:len(line[0])]
		name := line[1]
		org := line[2]
		consumption := line[3][0:strings.Index(line[3], "<")]

		fmt.Printf("%6s%16s%16s%16s\n",rank,name,org,consumption)
	}
	fmt.Println(separator)
}
