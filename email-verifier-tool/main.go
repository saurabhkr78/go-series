package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// Result holds all the DNS findings for a single domain.
// This struct is returned by checkDomain() and sent through the results channel.
type Result struct {
	Domain      string
	HasMX       bool
	HasSPF      bool
	SPFRecord   string
	HasDMARC    bool
	DMARCRecord string
}

// checkDomain performs three checks:
// 1. MX    → Does the domain have mail servers?
// 2. SPF   → Does the domain specify allowed mail senders?
// 3. DMARC → Does the domain enforce anti-spoofing policy?
// It then returns a Result struct.
func checkDomain(domain string) Result {
	var hasSPF, hasDMARC, hasMX bool
	var spfRecord, dmarcRecord string

	// --------------- MX CHECK ----------------
	// net.LookupMX queries DNS for MX (mail exchanger) records.
	// If MX records exist, domain can receive email.
	mxRecords, err := net.LookupMX(domain)
	if err == nil && len(mxRecords) > 0 {
		hasMX = true
	}

	// --------------- SPF CHECK ----------------
	// SPF is stored in DNS TXT records at the root domain.
	// We scan for a TXT record starting with "v=spf1".
	txtRecords, err := net.LookupTXT(domain)
	if err == nil {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf1") {
				hasSPF = true
				spfRecord = txt
				break
			}
		}
	}

	// --------------- DMARC CHECK ----------------
	// DMARC is always stored in: _dmarc.<domain>
	dmarcDomain := "_dmarc." + domain
	dmarcTxtRecords, err := net.LookupTXT(dmarcDomain)
	if err == nil {
		for _, txt := range dmarcTxtRecords {
			if strings.HasPrefix(txt, "v=DMARC1") {
				hasDMARC = true
				dmarcRecord = txt
				break
			}
		}
	}

	return Result{
		Domain:      domain,
		HasMX:       hasMX,
		HasSPF:      hasSPF,
		SPFRecord:   spfRecord,
		HasDMARC:    hasDMARC,
		DMARCRecord: dmarcRecord,
	}
}

// worker() is executed by each goroutine. It performs the following:
// 1. Continuously reads domains from the 'jobs' channel.
// 2. For each domain, runs checkDomain().
// 3. Sends the result into the 'results' channel.
//
// CHANNEL DIRECTIONS:
//
//	jobs <-chan string  → receive-only channel (worker can ONLY read)
//	results chan<- Result → send-only channel (worker can ONLY write)
//
// WORKFLOW:
//
//	jobs channel --> worker --> results channel
//
// When jobs channel closes, the loop ends and wg.Done() informs
// the WaitGroup that this worker has completed its work.
func worker(jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done() // tell the WaitGroup that THIS worker has finished

	for domain := range jobs { // This loop ends when jobs channel is closed
		results <- checkDomain(domain) // send result to output channel
	}
}

func main() {
	fmt.Println("domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord")

	scanner := bufio.NewScanner(os.Stdin)

	// -------------------------------------------------------------------------
	// CREATING CHANNELS
	//
	// jobs    → carries domains to be processed
	// results → carries results returned by workers
	//
	// BUFFERED CHANNELS (size = 100):
	// - Allow temporary storage without blocking.
	// - If buffer full: goroutine sleeps until space frees.
	// -------------------------------------------------------------------------
	jobs := make(chan string, 100)
	results := make(chan Result, 100)

	// -------------------------------------------------------------------------
	// WORKER POOL SETUP
	//
	// workerCount defines how many goroutines (workers) run in parallel.
	// More workers = faster processing, but too many may overwhelm DNS servers.
	//
	// WaitGroup ensures main() waits for ALL workers to finish.
	// -------------------------------------------------------------------------
	workerCount := 20
	var wg sync.WaitGroup

	// -------------------------------------------------------------------------
	// START WORKERS (concurrently)
	//
	// Each worker:
	// - Waits for a domain from jobs
	// - Processes it
	// - Sends a Result into results
	//
	// wg.Add(1) → register worker
	// go worker() → start a new goroutine
	// -------------------------------------------------------------------------
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// -------------------------------------------------------------------------
	// READING INPUT & SENDING JOBS
	//
	// This goroutine continuously reads domains from STDIN (scanner).
	// Each domain is pushed into jobs channel.
	//
	// When scanning ends, we close(jobs).
	// Closing jobs tells ALL workers: "No more domains are coming."
	// -------------------------------------------------------------------------
	go func() {
		for scanner.Scan() {
			jobs <- scanner.Text() // send domain to workers
		}
		close(jobs) // important: without this, workers never exit
	}()

	// -------------------------------------------------------------------------
	// CLOSING RESULTS CHANNEL
	//
	// This goroutine waits until ALL workers finish (wg.Wait()).
	// Only then it closes the results channel.
	//
	// Closing results signals the printing loop to stop.
	// -------------------------------------------------------------------------
	go func() {
		wg.Wait()
		close(results)
	}()

	// -------------------------------------------------------------------------
	// PRINT RESULTS
	//
	// Loop continues until results channel closes.
	// -------------------------------------------------------------------------
	for res := range results {
		fmt.Printf("%s,%t,%t,%s,%t,%s\n",
			res.Domain, res.HasMX, res.HasSPF, res.SPFRecord, res.HasDMARC, res.DMARCRecord)
	}

	// Check input errors
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}
