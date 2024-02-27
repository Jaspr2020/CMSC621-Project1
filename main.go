package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type job struct {
	datafilePathname string
	segmentStart     int
	segmentLength    int
}

type result struct {
	jobDescriptor job
	numPrimes     int
}

func main() {

	pathname := "C:\\Users\\Jacob\\go\\src\\CMSC621-Project1"
	M := 1
	N := 64000
	C := 1000

	readArgs(&pathname, &M, &N, &C)

	fmt.Println(pathname)
	fmt.Println(M)
	fmt.Println(N)
	fmt.Println(C)

	var wg sync.WaitGroup
	wg.Add(2 + M)

	// Create channels
	jobs := make(chan job)
	results := make(chan result)

	// Spawn one dispatcher thread
	go dispatcher(jobs)

	// Spawn M worker threads
	for i := 0; i < M; i++ {
		go worker(jobs, results)
	}

	// Spawn one consolidator thread
	go func() {
		go consolidator()
		wg.Done()
	}()

	wg.Wait()
}

func readArgs(pathname *string, M *int, N *int, C *int) {
	if len(os.Args) != 5 {
		return
	}

	// If all arguments are present, read them
	args := os.Args

	*pathname = args[1]

	// Check for valid M input
	i, err := strconv.Atoi(args[2])
	if err == nil {
		*M = i
	} else {
		os.Exit(1)
	}

	// Check for valid N input
	j, err := strconv.Atoi(args[3])
	if err == nil && j%8 == 0 {
		*N = j
	} else {
		os.Exit(1)
	}

	// Check for valid C input
	k, err := strconv.Atoi(args[4])
	if err == nil && k%8 == 0 {
		*C = k
	} else {
		os.Exit(1)
	}
}

func dispatcher(jobs chan<- job) {
	// Create job
	// Place job in queue
}

func worker(jobs <-chan job, results chan<- result) {
	// Sleep for between 400 and 600 msecs
	time.Sleep(time.Duration(rand.Intn(200)+400) * time.Millisecond)

	// Take a job from the jobqueue
	for j := range jobs {
		fmt.Println(j)
		// Read the job's segment from the datafile in chunks of C bytes (or less) at a time

		// Find the number of 64-bit unsigned integers in the segment that are prime numbers
		numPrimes := 0

		// Place a descriptor of it's finding/result to the results-queue
		results <- result{j, numPrimes}
	}
}

func consolidator() int {
	return 0
}
