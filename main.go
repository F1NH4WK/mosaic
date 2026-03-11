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

	"github.com/F1NH4WK/pupp/internal/cli"
	"github.com/F1NH4WK/pupp/internal/consumer"
	"github.com/F1NH4WK/pupp/internal/producer"
)

func main() {
	interactivePtr := flag.Bool("i", false, "Modo interativo (Perguntas passo a passo)")

	namePtr := flag.String("n", "", "Nome do alvo")
	yearPtr := flag.Int("y", 0, "Ano de nascimento")
	keywordsPtr := flag.String("k", "", "Palavras-chave (vírgula)")
	outputPtr := flag.String("o", "wordlist.txt", "Arquivo de saída")

	minLen := flag.Int("min", 6, "Tamanho mínimo da senha")
	reqUpper := flag.Bool("upper", false, "Exigir ao menos 1 letra maiúscula")
	reqLower := flag.Bool("lower", false, "Exigir ao menos 1 letra minúscula")
	reqNum := flag.Bool("num", false, "Exigir ao menos 1 número")
	reqSpec := flag.Bool("spec", false, "Exigir ao menos 1 caractere especial")
	verbosePtr := flag.Bool("v", false, "Modo verboso: Imprime as senhas geradas na tela (AVISO: Reduz drasticamente a velocidade)")


	flag.Parse()

	var rawWords []string
	var targetYear int

	if *interactivePtr {
		words, year := cli.StartInteractiveMode()
		if len(words) == 0 {
			fmt.Println("[-] Nenhuma informação fornecida. Abortando.")
			os.Exit(1)
		}
		rawWords = words
		targetYear = year
	} else {
		if *namePtr == "" {
			fmt.Println("Erro: O nome do alvo (-n) é obrigatório se não usar o modo interativo (-i).")
			os.Exit(1)
		}
		rawWords = append([]string{*namePtr}, strings.Split(*keywordsPtr, ",")...)
		targetYear = *yearPtr
	}

	rules := producer.Rules{
		MinLength:    *minLen,
		RequireUpper: *reqUpper,
		RequireLower: *reqLower,
		RequireNum:   *reqNum,
		RequireSpec:  *reqSpec,
	}

	start := time.Now()
	fmt.Println("[*] Gerando combinações base...")

	combinedBaseWords := producer.GenerateCombinations(rawWords, targetYear)

	passChan := make(chan string, 150000)
	var wg sync.WaitGroup

	doneChan := make(chan bool)
	go func() {
		if err := consumer.WriteToDisk(passChan, *outputPtr, *verbosePtr); err != nil {
			log.Fatalf("Erro ao escrever no disco: %v", err)
		}
		doneChan <- true
	}()

	jobs := make(chan string, len(combinedBaseWords))
	for _, w := range combinedBaseWords {
		jobs <- w
	}
	close(jobs)

	numWorkers := runtime.NumCPU()
	fmt.Printf("[*] Lançando %d Workers na CPU para aplicar Leetspeak...\n", numWorkers)

	for range numWorkers {
		wg.Go(func() {
			for baseWord := range jobs {
				producer.GeneratePasswords(baseWord, passChan, rules)
			}
		})
	}

	wg.Wait()
	close(passChan)
	<-doneChan

	fmt.Printf("[+] Wordlist pronta: %s | Tempo: %v\n", *outputPtr, time.Since(start))
}