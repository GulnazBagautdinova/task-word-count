package main

import (
	"os"
	"log"
	"sync"
	"fmt"
	"errors"
	"net/http"
	"regexp"
	"bufio"
	"io/ioutil"
)

type result struct{
	url string
	count int
	err error
}


func search (url string) result {

	var res result
	if url == "" {
		res.err = errors.New("empty url")
		return res
	}
	
	resp, err := http.Get(url)

	// Checks to know we have received a proper url
	if err != nil {
		res.err = errors.New("Cannot access site")
			res.url = url
			return res
		}
	
	defer resp.Body.Close()
  
	// Check the status code for a 200 to check we have received a proper response
	if resp.StatusCode != 200 {
		res.err = fmt.Errorf("HTTP Response Error %d\n", resp.StatusCode)
		res.url = url
		return res
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		// Check the page for the  term Go
		reg := regexp.MustCompile(`Go`)
		matched := reg.FindAllString(bodyString,-1)
		countmatched := len(matched)
		res = result{
			url: url,
			count: countmatched,
		}
	}
	return res
}

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Stdin: %s", err)
	}

	if (fi.Mode() & os.ModeCharDevice != 0){
		log.Fatal("data must be passed only from stdin")
	}
	results := make(chan result)
	doneChan := make(chan struct{})

	go func (results <- chan result, doneChan chan <- struct {}) {
		var total int
		for  res := range results{
			if res.err == nil{
			fmt.Printf ("Count for %s : %d\n" , res.url, res.count)
			} else {
		    fmt.Printf ("Count for %s : %s\n" , res.url,  res.err )
			}
			total +=  res.count
		}
		fmt.Printf("Total: %d\n", total)
		doneChan <- struct{}{}
	} (results, doneChan)

	// Setup a wait group in order to process all the urls
	var wg sync.WaitGroup
	var url string
	// Set the number of goroutines 
	// Processes urls
	goroutines := make(chan struct{}, 5)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url = scanner.Text()
		goroutines  <- struct{}{}
		wg.Add(1)
		go func (url string , results chan <- result, goroutines <- chan struct{}, wg *sync.WaitGroup) {
			results <- search (url)
			<- goroutines
			wg.Done()
		}(url , results, goroutines, &wg)
	}

	// Launch a goroutine to monitor when all the work is done
	// Wait for everything to be processed
	wg.Wait()
	close(goroutines)
	close(results)
	<- doneChan
}