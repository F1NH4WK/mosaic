package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/F1NH4WK/mosaic/internal/models"
)

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func StartInteractiveMode() (models.Profile, bool) {
	reader := bufio.NewReader(os.Stdin)
	var profile models.Profile

	fmt.Println("\n[+] Bem-vindo ao Mosaic (Modo Interativo).")
	fmt.Println("[+] Insira as informações sobre o alvo para gerar o dicionário.")
	fmt.Println("[+] Se não souber alguma informação, basta premir ENTER! ;)")

	if val := readLine(reader, "> Primeiro Nome: "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Sobrenome: "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Apelido (Nickname): "); val != "" {
		profile.Names = append(profile.Names, val)
	}

	dob := readLine(reader, "> Data de nascimento (DDMMAAAA): ")
	if dob != "" {
		cleanDOB := strings.ReplaceAll(dob, "-", "")
		cleanDOB = strings.ReplaceAll(cleanDOB, "/", "")
		profile.DOB = cleanDOB
	}

	fmt.Println()
	if val := readLine(reader, "> Nome do parceiro(a): "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Nome do filho(a): "); val != "" {
		profile.Names = append(profile.Names, val)
	}

	fmt.Println()
	if val := readLine(reader, "> Nome do Pet: "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Nome da Empresa: "); val != "" {
		profile.Names = append(profile.Names, val)
	}

	fmt.Println()
	if val := readLine(reader, "> Deseja adicionar palavras-chave extras? (separadas por vírgula): "); val != "" {
		extras := strings.Split(val, ",")
		for _, ex := range extras {
			profile.Keywords = append(profile.Keywords, strings.TrimSpace(ex))
		}
	}

	fmt.Println()
	useLeetspeak := false
	if val := readLine(reader, "> Deseja aplicar mutações Leetspeak pesadas? (s/N): "); strings.ToLower(val) == "s" {
		useLeetspeak = true
	}

	fmt.Println("\n[*] Inteligência recolhida com sucesso! A iniciar o gerador...")
	return profile, useLeetspeak
}