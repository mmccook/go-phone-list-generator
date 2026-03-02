package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
	"github.com/schollz/progressbar/v3"
)

func confirmSettings(areaCodes []string, count, countryCode int, fields []Field) (bool, error) {
	fmt.Println("\n=== Generation Settings ===")
	fmt.Printf("Area Codes: %s\n", strings.Join(areaCodes, ", "))
	fmt.Printf("Count: %d\n", count)
	fmt.Printf("Country Code: %d\n", countryCode)

	if len(fields) > 0 {
		fieldNames := make([]string, len(fields))
		for i, f := range fields {
			fieldNames[i] = f.Name
		}
		fmt.Printf("Additional Fields: %s\n", strings.Join(fieldNames, ", "))
	} else {
		fmt.Println("Additional Fields: None (phone numbers only)")
	}

	fmt.Print("\nProceed with generation? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("error reading input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

func generate(areaCodes []string, count, countryCode int, fields []Field, addrCache *addressCache) error {
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("phone_numbers_%s.csv", timestamp)

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Error closing file: %v\n", closeErr)
		}
	}()

	headers := []string{"ID", "phone Number"}
	for _, f := range fields {
		headers = append(headers, f.Name)
	}
	if _, err := file.WriteString(strings.Join(headers, ",") + "\n"); err != nil {
		return fmt.Errorf("error writing headers: %w", err)
	}

	log.Println("Running the generator...")

	bar := progressbar.NewOptions(count,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("[green][1/1][reset] Generating phone numbers..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Generate rows with even area code distribution
	areaCodeIndex := 0
	for i := 0; i < count; {
		// Use round-robin area code selection for even distribution
		selectedAreaCode := areaCodes[areaCodeIndex%len(areaCodes)]
		phone := RandomPhone(countryCode, []string{selectedAreaCode}, phonenumbers.E164)

		if phone.Possible() {
			row := []string{fmt.Sprintf("%d", i), phone.FormatedNumber()}

			// Add extra fields if requested
			if len(fields) > 0 {
				addrCache.reset()
				for _, f := range fields {
					row = append(row, f.Generator())
				}
			}

			if _, err := file.WriteString(strings.Join(row, ",") + "\n"); err != nil {
				return fmt.Errorf("error writing row: %w", err)
			}
			i++
			areaCodeIndex++
			_ = bar.Add(1)
		} else {
			// If phone is not possible, try next area code
			areaCodeIndex++
		}
	}

	_ = bar.Finish()
	fmt.Println() // Add newline after progress bar

	log.Printf("Done! Generated %d phone numbers in %s\n", count, fileName)
	return nil
}
