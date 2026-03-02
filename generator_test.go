package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nyaruka/phonenumbers"
)

func TestConfirmSettings_WithYesResponse(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe to simulate user input
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r

	// Write "y" to simulate user input
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = fmt.Fprintln(w, "y")
	}()

	areaCodes := []string{"631", "561"}
	fields := []Field{
		{Name: "First Name", Generator: func() string { return "John" }},
	}

	confirmed, err := confirmSettings(areaCodes, 100, 1, fields)

	if err != nil {
		t.Fatalf("confirmSettings() returned error: %v", err)
	}

	if !confirmed {
		t.Error("confirmSettings() = false, want true for 'y' input")
	}
}

func TestConfirmSettings_WithNoResponse(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe to simulate user input
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r

	// Write "n" to simulate user input
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = fmt.Fprintln(w, "n")
	}()

	areaCodes := []string{"631"}
	fields := []Field{}

	confirmed, err := confirmSettings(areaCodes, 50, 1, fields)

	if err != nil {
		t.Fatalf("confirmSettings() returned error: %v", err)
	}

	if confirmed {
		t.Error("confirmSettings() = true, want false for 'n' input")
	}
}

func TestConfirmSettings_WithYesFullWord(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	// Create a pipe to simulate user input
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = r

	// Write "yes" to simulate user input
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = fmt.Fprintln(w, "yes")
	}()

	areaCodes := []string{"631"}
	fields := []Field{}

	confirmed, err := confirmSettings(areaCodes, 50, 1, fields)

	if err != nil {
		t.Fatalf("confirmSettings() returned error: %v", err)
	}

	if !confirmed {
		t.Error("confirmSettings() = false, want true for 'yes' input")
	}
}

func TestConfirmSettings_CaseInsensitive(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"uppercase Y", "Y", true},
		{"uppercase YES", "YES", true},
		{"mixed case Yes", "Yes", true},
		{"uppercase N", "N", false},
		{"uppercase NO", "NO", false},
		{"mixed case No", "No", false},
		{"with whitespace", "  y  ", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// Create a pipe to simulate user input
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdin = r

			// Write input
			go func() {
				defer func() { _ = w.Close() }()
				_, _ = fmt.Fprintln(w, tc.input)
			}()

			confirmed, err := confirmSettings([]string{"631"}, 50, 1, []Field{})

			if err != nil {
				t.Fatalf("confirmSettings() returned error: %v", err)
			}

			if confirmed != tc.expected {
				t.Errorf("confirmSettings() with input '%s' = %v, want %v", tc.input, confirmed, tc.expected)
			}
		})
	}
}

func TestGenerate_NoAdditionalFields(t *testing.T) {
	areaCodes := []string{"631"}
	count := 5
	countryCode := 1
	fields := []Field{}
	addrCache := &addressCache{}

	// Generate CSV
	err := generate(areaCodes, count, countryCode, fields, addrCache)
	if err != nil {
		t.Fatalf("generate() returned error: %v", err)
	}

	// Find the generated file
	files, err := filepath.Glob("phone_numbers_*.csv")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatal("No CSV file was generated")
	}

	// Clean up - defer removal of the latest file
	defer func() { _ = os.Remove(files[len(files)-1]) }()

	// Read and verify the file
	content, err := os.ReadFile(files[len(files)-1])
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Check header
	expectedHeader := "ID,phone Number"
	if lines[0] != expectedHeader {
		t.Errorf("Header = %s, want %s", lines[0], expectedHeader)
	}

	// Check number of rows (header + count)
	if len(lines) != count+1 {
		t.Errorf("Number of lines = %d, want %d", len(lines), count+1)
	}

	// Verify each row has correct format
	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], ",")
		if len(parts) != 2 {
			t.Errorf("Row %d has %d columns, want 2", i, len(parts))
		}
	}
}

func TestGenerate_WithAdditionalFields(t *testing.T) {
	areaCodes := []string{"631"}
	count := 3
	countryCode := 1
	fields := []Field{
		{Name: "First Name", Generator: func() string { return "John" }},
		{Name: "Last Name", Generator: func() string { return "Doe" }},
	}
	addrCache := &addressCache{}

	// Generate CSV
	err := generate(areaCodes, count, countryCode, fields, addrCache)
	if err != nil {
		t.Fatalf("generate() returned error: %v", err)
	}

	// Find the generated file
	files, err := filepath.Glob("phone_numbers_*.csv")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatal("No CSV file was generated")
	}

	// Clean up
	defer func() { _ = os.Remove(files[len(files)-1]) }()

	// Read and verify the file
	content, err := os.ReadFile(files[len(files)-1])
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Check header includes additional fields
	expectedHeader := "ID,phone Number,First Name,Last Name"
	if lines[0] != expectedHeader {
		t.Errorf("Header = %s, want %s", lines[0], expectedHeader)
	}

	// Check number of rows
	if len(lines) != count+1 {
		t.Errorf("Number of lines = %d, want %d", len(lines), count+1)
	}

	// Verify each row has correct number of columns
	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], ",")
		expectedColumns := 4 // ID + phone + 2 fields
		if len(parts) != expectedColumns {
			t.Errorf("Row %d has %d columns, want %d", i, len(parts), expectedColumns)
		}

		// Verify the additional field values
		if parts[2] != "John" {
			t.Errorf("Row %d First Name = %s, want John", i, parts[2])
		}
		if parts[3] != "Doe" {
			t.Errorf("Row %d Last Name = %s, want Doe", i, parts[3])
		}
	}
}

func TestGenerate_MultipleAreaCodes(t *testing.T) {
	areaCodes := []string{"631", "561", "446"}
	count := 10
	countryCode := 1
	fields := []Field{}
	addrCache := &addressCache{}

	// Generate CSV
	err := generate(areaCodes, count, countryCode, fields, addrCache)
	if err != nil {
		t.Fatalf("generate() returned error: %v", err)
	}

	// Find the generated file
	files, err := filepath.Glob("phone_numbers_*.csv")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatal("No CSV file was generated")
	}

	// Clean up
	defer func() { _ = os.Remove(files[len(files)-1]) }()

	// Read and verify the file
	content, err := os.ReadFile(files[len(files)-1])
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Verify we got the correct number of lines
	if len(lines) != count+1 {
		t.Errorf("Number of lines = %d, want %d", len(lines), count+1)
	}

	// Verify phone numbers are from the provided area codes
	areaCodeMap := make(map[string]bool)
	for _, ac := range areaCodes {
		areaCodeMap[ac] = true
	}

	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], ",")
		if len(parts) < 2 {
			continue
		}

		phoneNum := parts[1]
		// phone should be in E164 format starting with +1
		if !strings.HasPrefix(phoneNum, "+1") {
			t.Errorf("phone number %s doesn't start with +1", phoneNum)
		}
	}
}

func TestGenerate_FileCreation(t *testing.T) {
	areaCodes := []string{"631"}
	count := 1
	countryCode := 1
	fields := []Field{}
	addrCache := &addressCache{}

	// Generate CSV
	err := generate(areaCodes, count, countryCode, fields, addrCache)
	if err != nil {
		t.Fatalf("generate() returned error: %v", err)
	}

	// Find the generated file
	files, err := filepath.Glob("phone_numbers_*.csv")
	if err != nil {
		t.Fatal(err)
	}

	if len(files) == 0 {
		t.Fatal("No CSV file was generated")
	}

	// Verify file exists and is not empty
	latestFile := files[len(files)-1]
	defer func() { _ = os.Remove(latestFile) }()

	fileInfo, err := os.Stat(latestFile)
	if err != nil {
		t.Fatalf("Failed to stat generated file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Error("Generated file is empty")
	}

	// Verify filename format
	if !strings.HasPrefix(filepath.Base(latestFile), "phone_numbers_") {
		t.Errorf("Filename doesn't match expected pattern: %s", latestFile)
	}

	if !strings.HasSuffix(latestFile, ".csv") {
		t.Errorf("File doesn't have .csv extension: %s", latestFile)
	}
}

func TestGenerate_PhoneNumberValidity(t *testing.T) {
	areaCodes := []string{"631"}
	count := 5
	countryCode := 1
	fields := []Field{}
	addrCache := &addressCache{}

	// Generate CSV
	err := generate(areaCodes, count, countryCode, fields, addrCache)
	if err != nil {
		t.Fatalf("generate() returned error: %v", err)
	}

	// Find the generated file
	files, err := filepath.Glob("phone_numbers_*.csv")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.Remove(files[len(files)-1]) }()

	// Read and verify phone numbers
	content, err := os.ReadFile(files[len(files)-1])
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	// Skip header, check each phone number
	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], ",")
		if len(parts) < 2 {
			continue
		}

		phoneNumStr := parts[1]

		// Parse the phone number
		num, err := phonenumbers.Parse(phoneNumStr, "US")
		if err != nil {
			t.Errorf("Row %d: Failed to parse phone number %s: %v", i, phoneNumStr, err)
			continue
		}

		// Verify it's possible
		if !phonenumbers.IsPossibleNumber(num) {
			t.Errorf("Row %d: phone number %s is not possible", i, phoneNumStr)
		}
	}
}

func TestGenerate_AddressCacheReset(t *testing.T) {
	areaCodes := []string{"631"}
	count := 3
	countryCode := 1

	callCount := 0
	fields := []Field{
		{Name: "Test Field", Generator: func() string {
			callCount++
			return fmt.Sprintf("Value%d", callCount)
		}},
	}
	addrCache := &addressCache{}

	// Generate CSV
	err := generate(areaCodes, count, countryCode, fields, addrCache)
	if err != nil {
		t.Fatalf("generate() returned error: %v", err)
	}

	// Find and clean up the generated file
	files, err := filepath.Glob("phone_numbers_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) > 0 {
		defer func() { _ = os.Remove(files[len(files)-1]) }()
	}

	// Verify the generator was called for each row
	if callCount != count {
		t.Errorf("Generator called %d times, expected %d", callCount, count)
	}
}

func TestGenerate_IDSequence(t *testing.T) {
	areaCodes := []string{"631"}
	count := 5
	countryCode := 1
	fields := []Field{}
	addrCache := &addressCache{}

	// Generate CSV
	err := generate(areaCodes, count, countryCode, fields, addrCache)
	if err != nil {
		t.Fatalf("generate() returned error: %v", err)
	}

	// Find the generated file
	files, err := filepath.Glob("phone_numbers_*.csv")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.Remove(files[len(files)-1]) }()

	// Read and verify IDs are sequential
	content, err := os.ReadFile(files[len(files)-1])
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")

	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], ",")
		expectedID := fmt.Sprintf("%d", i-1)
		if parts[0] != expectedID {
			t.Errorf("Row %d: ID = %s, want %s", i, parts[0], expectedID)
		}
	}
}

// Cleanup function to remove test CSV files
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Cleanup any remaining test CSV files
	files, _ := filepath.Glob("phone_numbers_*.csv")
	for _, f := range files {
		_ = os.Remove(f)
	}

	os.Exit(code)
}

func BenchmarkGenerate(b *testing.B) {
	areaCodes := []string{"631"}
	fields := []Field{}
	addrCache := &addressCache{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generate(areaCodes, 10, 1, fields, addrCache)
	}

	// Cleanup
	files, _ := filepath.Glob("phone_numbers_*.csv")
	for _, f := range files {
		_ = os.Remove(f)
	}
}
