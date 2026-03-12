package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/F1NH4WK/mosaic/internal/cli"
	"github.com/F1NH4WK/mosaic/internal/consumer"
	"github.com/F1NH4WK/mosaic/internal/models"
	"github.com/F1NH4WK/mosaic/internal/producer"
)

func main() {
	interactivePtr := flag.Bool("i", false, "Interactive mode (Step-by-step prompts)")

	namePtr := flag.String("n", "", "Target's first name")
	surnamePtr := flag.String("s", "", "Target's surname")
	nickPtr := flag.String("nick", "", "Target's nickname")
	dobPtr := flag.String("dob", "", "Date of birth (DDMMYYYY)")
	partnerPtr := flag.String("p", "", "Partner's name")
	childPtr := flag.String("c", "", "Child's name")
	keywordsPtr := flag.String("k", "", "Extra keywords separated by comma")
	
	outputPtr := flag.String("o", "wordlist.txt", "Output file")
	leetPtr := flag.Bool("leet", false, "Enable Leetspeak mutations (Warning: Generates massive lists)")
	verbosePtr := flag.Bool("v", false, "Verbose mode: Print passwords to screen (Warning: Reduces I/O speed)")

	minLen := flag.Int("min", 6, "Minimum password length")
	reqUpper := flag.Bool("upper", false, "Require at least 1 uppercase letter")
	reqLower := flag.Bool("lower", false, "Require at least 1 lowercase letter")
	reqNum := flag.Bool("num", false, "Require at least 1 number")
	reqSpec := flag.Bool("spec", false, "Require at least 1 special character")

	flag.Parse()

	var targetProfile models.Profile
	var applyLeetspeak bool

	if *interactivePtr {
		profile, leet := cli.StartInteractiveMode()
		if len(profile.Names) == 0 && len(profile.Keywords) == 0 {
			fmt.Println("[-] No information provided. Aborting.")
			os.Exit(1)
		}
		targetProfile = profile
		applyLeetspeak = leet
	} else {
		if *namePtr == "" {
			fmt.Println("Error: Target's first name (-n) is required outside interactive mode (-i).")
			os.Exit(1)
		}

		var names []string
		if *namePtr != "" { names = append(names, *namePtr) }
		if *surnamePtr != "" { names = append(names, *surnamePtr) }
		if *nickPtr != "" { names = append(names, *nickPtr) }
		if *partnerPtr != "" { names = append(names, *partnerPtr) }
		if *childPtr != "" { names = append(names, *childPtr) }

		var keywords []string
		if *keywordsPtr != "" {
			for _, k := range strings.Split(*keywordsPtr, ",") {
				keywords = append(keywords, strings.TrimSpace(k))
			}
		}

		cleanDOB := strings.ReplaceAll(*dobPtr, "-", "")
		cleanDOB = strings.ReplaceAll(cleanDOB, "/", "")

		targetProfile = models.Profile{
			Names:    names,
			DOB:      cleanDOB,
			Keywords: keywords,
		}
		applyLeetspeak = *leetPtr
	}

	rules := producer.Rules{
		MinLength:    *minLen,
		RequireUpper: *reqUpper,
		RequireLower: *reqLower,
		RequireNum:   *reqNum,
		RequireSpec:  *reqSpec,
		UseLeetspeak: applyLeetspeak,
	}

	start := time.Now()
	fmt.Println("[*] Starting Mosaic - Tactical Profiler & Generator...")
	fmt.Println("[*] Analyzing profile and extracting base heuristics...")

	combinedBaseWords := producer.GenerateCombinations(targetProfile)
	fmt.Printf("[*] Generated %d precise base permutations.\n", len(combinedBaseWords))

	passChan := make(chan string, 150000)
	var wg sync.WaitGroup

	doneChan := make(chan bool)
	go func() {
		if err := consumer.WriteToDisk(passChan, *outputPtr, *verbosePtr); err != nil {
			log.Fatalf("Error writing to disk: %v", err)
		}
		doneChan <- true
	}()

	jobs := make(chan string, len(combinedBaseWords))
	for _, w := range combinedBaseWords {
		jobs <- w
	}
	close(jobs)

	numWorkers := runtime.NumCPU()
	if applyLeetspeak {
		fmt.Printf("[*] Leetspeak Mode ON: Launching %d CPU Workers for deep mutation...\n", numWorkers)
	} else {
		fmt.Printf("[*] Leetspeak Mode OFF: Writing wordlist directly...\n")
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for baseWord := range jobs {
				producer.GeneratePasswords(baseWord, passChan, rules)
			}
		}()
	}

	wg.Wait()
	close(passChan)
	<-doneChan

	fmt.Printf("[+] Wordlist ready: %s | Execution Time: %v\n", *outputPtr, time.Since(start))
}