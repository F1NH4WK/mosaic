// main.go
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
	interactivePtr := flag.Bool("i", false, "Modo interativo (Perguntas passo a passo)")

	namePtr := flag.String("n", "", "Primeiro nome do alvo")
	surnamePtr := flag.String("s", "", "Sobrenome do alvo")
	nickPtr := flag.String("nick", "", "Apelido (Nickname)")
	dobPtr := flag.String("dob", "", "Data de nascimento (DDMMAAAA)")
	partnerPtr := flag.String("p", "", "Nome do parceiro(a)")
	childPtr := flag.String("c", "", "Nome do filho(a)")
	keywordsPtr := flag.String("k", "", "Palavras-chave extras separadas por vírgula")
	
	outputPtr := flag.String("o", "wordlist.txt", "Ficheiro de saída")
	leetPtr := flag.Bool("leet", false, "Ativar mutações de Leetspeak (Gera listas massivas)")
	verbosePtr := flag.Bool("v", false, "Modo verboso: Imprime as senhas na tela (Aviso: Reduz a velocidade)")

	minLen := flag.Int("min", 6, "Tamanho mínimo da palavra-passe")
	reqUpper := flag.Bool("upper", false, "Exigir pelo menos 1 letra maiúscula")
	reqLower := flag.Bool("lower", false, "Exigir pelo menos 1 letra minúscula")
	reqNum := flag.Bool("num", false, "Exigir pelo menos 1 número")
	reqSpec := flag.Bool("spec", false, "Exigir pelo menos 1 caractere especial")

	flag.Parse()

	var targetProfile models.Profile
	var applyLeetspeak bool

	if *interactivePtr {
		profile, leet := cli.StartInteractiveMode()
		if len(profile.Names) == 0 && len(profile.Keywords) == 0 {
			fmt.Println("[-] Nenhuma informação fornecida. A abortar.")
			os.Exit(1)
		}
		targetProfile = profile
		applyLeetspeak = leet
	} else {
		if *namePtr == "" {
			fmt.Println("Erro: O primeiro nome do alvo (-n) é obrigatório fora do modo interativo (-i).")
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
	fmt.Println("[*] A iniciar o Mosaic - Gerador e Perfilador Tático...")
	fmt.Println("[*] A analisar perfil e a extrair heurísticas base...")

	combinedBaseWords := producer.GenerateCombinations(targetProfile)
	fmt.Printf("[*] Foram geradas %d permutações base precisas.\n", len(combinedBaseWords))

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
	if applyLeetspeak {
		fmt.Printf("[*] Modo Leetspeak ON: A lançar %d Workers na CPU para mutação profunda...\n", numWorkers)
	} else {
		fmt.Printf("[*] Modo Leetspeak OFF: A gravar dicionário diretamente...\n")
	}

	for range numWorkers {
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

	fmt.Printf("[+] Dicionário pronto: %s | Tempo de Execução: %v\n", *outputPtr, time.Since(start))
}