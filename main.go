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
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		// Check the page for the search term.
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
	var f *os.File	
	f = os.Stdin
	defer f.Close()
	//var urls []string		
	results := make(chan result)
	doneChan := make(chan struct{})
    
	  go func (results <- chan result, doneChan chan <- struct {}) {
		var total int
         for  res := range results{
			 fmt.Printf ("Count for %s  - %d\n" , res.url, res.count )
			 total +=  res.count
		 }    
		 fmt.Printf("Total count: %d\n", total)
		 doneChan <- struct{}{}
	  } (results, doneChan)

	// Setup a wait group so we can process all the feeds.
	var wg sync.WaitGroup
	var url string	
	
	goroutines := make(chan struct{}, 5)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {	
	url = scanner.Text()
  // Set the number of goroutines we need to wait for while
<<<<<<< HEAD
<<<<<<< HEAD
    goroutines  <- struct{}{}
=======
       goroutines  <- struct{}{}
>>>>>>> f9df75b2c14a9fa62bd88e551a26d4cc20738d29
=======
       goroutines  <- struct{}{}
>>>>>>> f9df75b2c14a9fa62bd88e551a26d4cc20738d29
       wg.Add(1) 
     go func (url string , results chan <- result, goroutines <- chan struct{}, wg *sync.WaitGroup) {
	    results <- search (url)   
	    <- goroutines
	    wg.Done() 
       }(url , results, goroutines, &wg)
	 }

// Launch a goroutine to monitor when all the work is done.
	// Wait for everything to be processed.
	wg.Wait()
	close(goroutines)
	close(results)
    <- doneChan
<<<<<<< HEAD
<<<<<<< HEAD
}
=======
}
>>>>>>> f9df75b2c14a9fa62bd88e551a26d4cc20738d29
=======
}
>>>>>>> f9df75b2c14a9fa62bd88e551a26d4cc20738d29
