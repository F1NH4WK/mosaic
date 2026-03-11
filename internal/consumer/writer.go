package consumer

import (
	"bufio"
	"fmt"
	"os"
)

func WriteToDisk(inChan <-chan string, filename string, verbose bool) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for password := range inChan {
		writer.WriteString(password + "\n")
		

		if verbose {
			fmt.Println(password)
		}
	}

	return nil
}