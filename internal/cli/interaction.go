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

	fmt.Println("\n[+] Welcome to Mosaic (Interactive Mode).")
	fmt.Println("[+] Insert the target's information to generate the wordlist.")
	fmt.Println("[+] If you don't know some info, just press ENTER! ;)")

	if val := readLine(reader, "> First Name: "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Surname: "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Nickname: "); val != "" {
		profile.Names = append(profile.Names, val)
	}

	dob := readLine(reader, "> Date of birth (DDMMYYYY): ")
	if dob != "" {
		cleanDOB := strings.ReplaceAll(dob, "-", "")
		cleanDOB = strings.ReplaceAll(cleanDOB, "/", "")
		profile.DOB = cleanDOB
	}

	fmt.Println()
	if val := readLine(reader, "> Partner's name: "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Child's name: "); val != "" {
		profile.Names = append(profile.Names, val)
	}

	fmt.Println()
	if val := readLine(reader, "> Pet's name: "); val != "" {
		profile.Names = append(profile.Names, val)
	}
	if val := readLine(reader, "> Company name: "); val != "" {
		profile.Names = append(profile.Names, val)
	}

	fmt.Println()
	if val := readLine(reader, "> Do you want to add extra keywords? (comma-separated): "); val != "" {
		extras := strings.Split(val, ",")
		for _, ex := range extras {
			profile.Keywords = append(profile.Keywords, strings.TrimSpace(ex))
		}
	}

	fmt.Println()
	useLeetspeak := false
	if val := readLine(reader, "> Do you want to apply heavy Leetspeak mutations? (y/N): "); strings.ToLower(val) == "y" {
		useLeetspeak = true
	}

	fmt.Println("\n[*] Intelligence gathered successfully! Starting the generator...")
	return profile, useLeetspeak
}