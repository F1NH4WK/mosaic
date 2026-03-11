package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func StartInteractiveMode() ([]string, int, bool) {
	reader := bufio.NewReader(os.Stdin)
	var words []string
	var year int

	fmt.Println("\n[+] Insira as informações sobre o alvo para gerar o dicionário.")
	fmt.Println("[+] Se não souber alguma informação, basta apertar ENTER! ;)")

	if val := readLine(reader, "> Primeiro Nome: "); val != "" {
		words = append(words, val)
	}
	if val := readLine(reader, "> Sobrenome: "); val != "" {
		words = append(words, val)
	}
	if val := readLine(reader, "> Apelido (Nickname): "); val != "" {
		words = append(words, val)
	}

	birthdate := readLine(reader, "> Data de nascimento (DDMMAAAA): ")
	if len(birthdate) == 8 {
		if y, err := strconv.Atoi(birthdate[4:]); err == nil {
			year = y
		}
		words = append(words, birthdate[:2], birthdate[2:4], birthdate[4:])
	}

	fmt.Println()
	if val := readLine(reader, "> Nome do parceiro(a): "); val != "" {
		words = append(words, val)
	}
	if val := readLine(reader, "> Nome do filho(a): "); val != "" {
		words = append(words, val)
	}

	fmt.Println()
	if val := readLine(reader, "> Nome do Pet: "); val != "" {
		words = append(words, val)
	}
	if val := readLine(reader, "> Nome da Empresa: "); val != "" {
		words = append(words, val)
	}

	fmt.Println()
	if val := readLine(reader, "> Deseja adicionar palavras-chave extras? (separadas por vírgula): "); val != "" {
		extras := strings.Split(val, ",")
		for _, ex := range extras {
			words = append(words, strings.TrimSpace(ex))
		}
	}

	fmt.Println()
	useLeetspeak := false
	if val := readLine(reader, "> Deseja aplicar mutações Leetspeak pesadas? (s/N): "); strings.ToLower(val) == "s" {
		useLeetspeak = true
	}

	fmt.Println("\n[*] Informações coletadas com sucesso!")
	return words, year, useLeetspeak
}