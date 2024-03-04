package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type job struct {
	datafilePathname string
	segmentStart     int64
	segmentLength    int64
}

type result struct {
	jobDescriptor job
	numPrimes     int
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	defer timer("main")()

	pathname := "C:\\Users\\Jacob\\go\\src\\CMSC621-Project1\\filename.dat"
	M := 10000
	N := int64(64000 * 1)
	C := int64(1000 * 1)

	readArgs(&pathname, &M, &N, &C)

	var wg sync.WaitGroup

	// Create channels
	jobs := make(chan job)
	results := make(chan result)
	totals := make(chan int, 2)

	// Spawn one dispatcher thread
	go dispatcher(jobs, pathname, N)

	// Spawn M worker threads
	for i := 0; i < M; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(jobs, results, C)
		}()
	}

	// Spawn one consolidator thread
	go consolidator(results, totals)

	wg.Wait()
	close(results)

	fmt.Println("Total primes:", <-totals)
	fmt.Println("Number of jobs:", <-totals)
}

func readArgs(pathname *string, M *int, N *int64, C *int64) {
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
		*N = int64(j)
	} else {
		os.Exit(1)
	}

	// Check for valid C input
	k, err := strconv.Atoi(args[4])
	if err == nil && k%8 == 0 {
		*C = int64(k)
	} else {
		os.Exit(1)
	}
}

func dispatcher(jobs chan<- job, pathname string, N int64) {
	// Read file
	file, err := os.Open(pathname)
	if err != nil {
		os.Exit(1)
	}

	stats, statsErr := file.Stat()
	if statsErr != nil {
		os.Exit(1)
	}

	var size int64 = stats.Size()

	var i int64 = 0
	for i < size {
		// Create job
		j := job{datafilePathname: pathname, segmentStart: i, segmentLength: min(N, size-i)}
		// Place job in queue
		jobs <- j
		i += N
	}
	close(jobs)
}

func worker(jobs <-chan job, results chan<- result, C int64) {
	// Sleep for between 400 and 600 msecs
	time.Sleep(time.Duration(rand.Intn(200)+400) * time.Millisecond)

	// Take a job from the jobqueue
	for j := range jobs {
		// Primes counter
		numPrimes := 0

		// Read the job's segment from the datafile in chunks of C bytes (or less) at a time
		i := int64(0)
		for i < j.segmentLength {
			file, err := os.Open(j.datafilePathname)
			if err != nil {
				os.Exit(1)
			}
			buffer := make([]byte, min(C, j.segmentLength-i))
			length, err := file.ReadAt(buffer, j.segmentStart+i)
			if err != nil {
				os.Exit(1)
			}

			// Find the number of 64-bit unsigned integers in the segment that are prime numbers
			k := 0
			for k < length {
				var result uint64
				err = binary.Read(bytes.NewReader(buffer[k:min(k+8, length)]), binary.LittleEndian, &result)
				if err != nil {
					os.Exit(1)
				}
				k += 8

				var num big.Int
				num.SetUint64(result)
				if num.ProbablyPrime(0) {
					numPrimes++
				}
			}

			i += C
		}

		// Place a descriptor of it's finding/result to the results-queue
		results <- result{j, numPrimes}
	}
}

func consolidator(results <-chan result, totals chan<- int) {
	// Take a result from the results queue
	numPrimes := 0
	numResults := 0
	for r := range results {
		numPrimes += r.numPrimes
		numResults++
		fmt.Println("Running Total:", numPrimes)
	}

	totals <- numPrimes
	totals <- numResults
}
