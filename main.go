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
	"github.com/F1NH4WK/mosaic/internal/producer"
)

func main() {
	interactivePtr := flag.Bool("i", false, "Modo interativo (Perguntas passo a passo)")

	namePtr := flag.String("n", "", "Nome do alvo")
	yearPtr := flag.Int("y", 0, "Ano de nascimento")
	keywordsPtr := flag.String("k", "", "Palavras-chave (vírgula)")
	outputPtr := flag.String("o", "passwords.txt", "Ficheiro de saída")
	leetPtr := flag.Bool("leet", false, "Ativar mutações de Leetspeak (Atenção: gera dicionários massivos)")

	// Restrictions flags
	minLen := flag.Int("min", 6, "Tamanho mínimo da senha")
	reqUpper := flag.Bool("upper", false, "Exigir ao menos 1 letra maiúscula")
	reqLower := flag.Bool("lower", false, "Exigir ao menos 1 letra minúscula")
	reqNum := flag.Bool("num", false, "Exigir ao menos 1 número")
	reqSpec := flag.Bool("spec", false, "Exigir ao menos 1 caractere especial")

	flag.Parse()

	var rawWords []string
	var targetYear int
	var applyLeetspeak bool

	if *interactivePtr {
		words, year, leet := cli.StartInteractiveMode()
		if len(words) == 0 {
			fmt.Println("[-] Nenhuma informação fornecida. A abortar.")
			os.Exit(1)
		}
		rawWords = words
		targetYear = year
		applyLeetspeak = leet
	} else {
		if *namePtr == "" {
			fmt.Println("Erro: O nome do alvo (-n) é obrigatório se não usar o modo interativo (-i).")
			os.Exit(1)
		}
		rawWords = append([]string{*namePtr}, strings.Split(*keywordsPtr, ",")...)
		targetYear = *yearPtr
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
	fmt.Println("[*] A iniciar o PersonaForge - O gerador tático...")
	fmt.Println("[*] A cruzar palavras e a gerar combinações base...")

	combinedBaseWords := producer.GenerateCombinations(rawWords, targetYear)

	passChan := make(chan string, 150000)
	var wg sync.WaitGroup

	doneChan := make(chan bool)
	go func() {
		if err := consumer.WriteToDisk(passChan, *outputPtr, false); err != nil {
			log.Fatalf("Erro: %v", err)
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
		fmt.Printf("[*] Modo Leetspeak ON: A lançar %d Workers na CPU...\n", numWorkers)
	} else {
		fmt.Printf("[*] Modo Leetspeak OFF: A extrair combinações diretas...\n")
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

	fmt.Printf("[+] Dicionário pronto: %s | Tempo: %v\n", *outputPtr, time.Since(start))
}