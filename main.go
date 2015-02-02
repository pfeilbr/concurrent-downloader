package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type URLJob struct {
	url string
}

type URLResult struct {
	URL  string
	Body string
	MD5  string
}

type TaskRunner struct {
	parallelLimit int
}

func fetch(url string) URLResult {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	hash := md5.Sum(body)
	bodyString := string(body)[0:10] + " ... truncated for readability ..."
	r := URLResult{URL: url, Body: bodyString, MD5: fmt.Sprintf("%x", hash)}
	//fmt.Println(r)
	return r
}

func worker(id int, jobs <-chan URLJob, results chan<- URLResult) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j.url)
		results <- fetch(j.url)
	}
}

func (tr TaskRunner) Run() {

	totalJobs := 50
	parallelLimit := tr.parallelLimit
	jobs := make(chan URLJob, parallelLimit)
	results := make(chan URLResult, totalJobs)
	urlResultList := make([]URLResult, 0)

	// create the jobs
	var urlJobList = make([]URLJob, 0)
	for i := 1; i <= totalJobs; i++ {
		url := fmt.Sprintf("http://www.google.com/?q=%d", rand.Intn(1000))
		urlJobList = append(urlJobList, URLJob{url: url})
	}

	// start <parallelLimit> # of workers
	// they will be blocked waiting for work at this point
	for w := 1; w <= parallelLimit; w++ {
		go worker(w, jobs, results)
	}

	// start sending jobs to be picked up by our workers
	for _, urlJob := range urlJobList {
		jobs <- urlJob
	}

	for i := 1; i < totalJobs; i++ {
		select {
		case urlResult := <-results:
			urlResultList = append(urlResultList, urlResult)
		}
	}

	// convert array to JSON.  MarshalIndent does JSON pretty print
	jsonBytes, err := json.MarshalIndent(urlResultList, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}

func main() {

	TaskRunner{parallelLimit: 10}.Run()

}
