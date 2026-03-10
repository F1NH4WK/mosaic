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

	"github.com/F1NH4WK/pupp/internal/consumer"
	"github.com/F1NH4WK/pupp/internal/producer"
)

func main() {
	// Flags antigas
	namePtr := flag.String("n", "", "Nome do alvo")
	yearPtr := flag.Int("y", 0, "Ano de nascimento")
	keywordsPtr := flag.String("k", "", "Palavras-chave (vírgula)")
	outputPtr := flag.String("o", "wordlist.txt", "Arquivo de saída")

	// NOVAS FLAGS DE RESTRIÇÃO
	minLen := flag.Int("min", 6, "Tamanho mínimo da senha")
	reqUpper := flag.Bool("upper", false, "Exigir ao menos 1 letra maiúscula")
	reqLower := flag.Bool("lower", false, "Exigir ao menos 1 letra minúscula")
	reqNum := flag.Bool("num", false, "Exigir ao menos 1 número")
	reqSpec := flag.Bool("spec", false, "Exigir ao menos 1 caractere especial")

	flag.Parse()

	if *namePtr == "" {
		fmt.Println("Erro: O nome do alvo (-n) é obrigatório.")
		os.Exit(1)
	}

	rules := producer.Rules{
		MinLength:    *minLen,
		RequireUpper: *reqUpper,
		RequireLower: *reqLower,
		RequireNum:   *reqNum,
		RequireSpec:  *reqSpec,
	}

	var keywords []string
	if *keywordsPtr != "" {
		keywords = strings.Split(*keywordsPtr, ",")
	}

	start := time.Now()
	fmt.Println("[*] Gerando combinações base...")

	rawWords := append([]string{*namePtr}, keywords...)
	combinedBaseWords := producer.GenerateCombinations(rawWords, *yearPtr)

	// Otimizando o buffer baseado na matemática: 150k é um excelente sweet-spot.
	passChan := make(chan string, 150000)
	var wg sync.WaitGroup

	doneChan := make(chan bool)
	go func() {
		if err := consumer.WriteToDisk(passChan, *outputPtr); err != nil {
			log.Fatalf("Erro: %v", err)
		}
		doneChan <- true
	}()

	// IMPLEMENTAÇÃO DO WORKER POOL
	// 1. Canal de tarefas contendo as palavras base
	jobs := make(chan string, len(combinedBaseWords))
	for _, w := range combinedBaseWords {
		jobs <- w
	}
	close(jobs) // Fecha o canal de jobs indicando que não há novas bases

	// 2. Descobre quantos núcleos lógicos a máquina tem
	numWorkers := runtime.NumCPU()
	fmt.Printf("[*] Lançando %d Workers na CPU para aplicar Leetspeak...\n", numWorkers)

	// 3. Inicia os Workers limitados
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for baseWord := range jobs {
				// Os workers pegam uma palavra da fila, processam e pedem a próxima
				producer.GeneratePasswords(baseWord, passChan, rules)
			}
		}()
	}

	wg.Wait()
	close(passChan)
	<-doneChan

	fmt.Printf("[+] Wordlist pronta: %s | Tempo: %v\n", *outputPtr, time.Since(start))
}