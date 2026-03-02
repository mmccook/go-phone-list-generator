package main

import (
	"flag"
	"log"
	"strings"
)

func main() {
	areacodes := flag.String("areacodes", "631,561,446", "Area codes to generate phone numbers from")
	count := flag.Int("count", 25000, "Number of phone numbers to generate")
	countryCode := flag.Int("countryCode", 1, "Country Code to use, default is 1 for US")
	fields := flag.String("fields", "", "Comma-separated list of fields to include (e.g., 'firstname,lastname,email'). Leave empty for phone numbers only. Available: firstname, lastname, address, city, state, zip, latitude, longitude, email, dob, username, company")

	flag.Parse()

	addrCache := &addressCache{}
	fieldRegistry := getFieldRegistry(addrCache)

	selectedFields, err := parseRequestedFields(*fields, fieldRegistry)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	areaCodesToUse := strings.Split(*areacodes, ",")

	confirmed, err := confirmSettings(areaCodesToUse, *count, *countryCode, selectedFields)
	if err != nil {
		log.Fatalf("Error reading confirmation: %v\n", err)
	}

	if !confirmed {
		log.Println("Generation cancelled by user.")
		return
	}

	if err := generate(areaCodesToUse, *count, *countryCode, selectedFields, addrCache); err != nil {
		log.Fatalf("Error generating CSV: %v\n", err)
	}
}
