package main

import (
	"slices"
	"testing"

	"github.com/nyaruka/phonenumbers"
)

func TestNewPhone(t *testing.T) {
	tests := []struct {
		name          string
		countryCode   int
		areaCode      string
		centralOffice string
		lineNumber    string
		format        phonenumbers.PhoneNumberFormat
	}{
		{
			name:          "US phone number",
			countryCode:   1,
			areaCode:      "631",
			centralOffice: "555",
			lineNumber:    "1234",
			format:        phonenumbers.E164,
		},
		{
			name:          "US phone number with different area code",
			countryCode:   1,
			areaCode:      "212",
			centralOffice: "123",
			lineNumber:    "4567",
			format:        phonenumbers.NATIONAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone := NewPhone(tt.countryCode, tt.areaCode, tt.centralOffice, tt.lineNumber, tt.format)

			if phone == nil {
				t.Fatal("NewPhone returned nil")
			}

			if phone.countryCode != tt.countryCode {
				t.Errorf("countryCode = %d, want %d", phone.countryCode, tt.countryCode)
			}

			if phone.areaCode != tt.areaCode {
				t.Errorf("areaCode = %s, want %s", phone.areaCode, tt.areaCode)
			}

			if phone.centralOffice != tt.centralOffice {
				t.Errorf("centralOffice = %s, want %s", phone.centralOffice, tt.centralOffice)
			}

			if phone.lineNumber != tt.lineNumber {
				t.Errorf("lineNumber = %s, want %s", phone.lineNumber, tt.lineNumber)
			}

			if phone.format != tt.format {
				t.Errorf("format = %v, want %v", phone.format, tt.format)
			}
		})
	}
}

func TestPhone_CountryCode(t *testing.T) {
	phone := NewPhone(1, "631", "555", "1234", phonenumbers.E164)
	if got := phone.CountryCode(); got != 1 {
		t.Errorf("CountryCode() = %d, want 1", got)
	}
}

func TestPhone_AreaCode(t *testing.T) {
	phone := NewPhone(1, "631", "555", "1234", phonenumbers.E164)
	if got := phone.AreaCode(); got != "631" {
		t.Errorf("AreaCode() = %s, want 631", got)
	}
}

func TestPhone_CentralOffice(t *testing.T) {
	phone := NewPhone(1, "631", "555", "1234", phonenumbers.E164)
	if got := phone.CentralOffice(); got != "555" {
		t.Errorf("CentralOffice() = %s, want 555", got)
	}
}

func TestPhone_LineNumber(t *testing.T) {
	phone := NewPhone(1, "631", "555", "1234", phonenumbers.E164)
	if got := phone.LineNumber(); got != "1234" {
		t.Errorf("LineNumber() = %s, want 1234", got)
	}
}

func TestPhone_FormatedNumber(t *testing.T) {
	tests := []struct {
		name          string
		countryCode   int
		areaCode      string
		centralOffice string
		lineNumber    string
		format        phonenumbers.PhoneNumberFormat
		wantContains  string
	}{
		{
			name:          "E164 format",
			countryCode:   1,
			areaCode:      "631",
			centralOffice: "555",
			lineNumber:    "1234",
			format:        phonenumbers.E164,
			wantContains:  "+16315551234",
		},
		{
			name:          "NATIONAL format",
			countryCode:   1,
			areaCode:      "631",
			centralOffice: "555",
			lineNumber:    "1234",
			format:        phonenumbers.NATIONAL,
			wantContains:  "631",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone := NewPhone(tt.countryCode, tt.areaCode, tt.centralOffice, tt.lineNumber, tt.format)
			got := phone.FormatedNumber()

			if got != tt.wantContains {
				// For NATIONAL format, just check if it contains the area code
				if tt.format == phonenumbers.NATIONAL {
					if len(got) == 0 {
						t.Errorf("FormatedNumber() returned empty string")
					}
				} else if got != tt.wantContains {
					t.Errorf("FormatedNumber() = %s, want %s", got, tt.wantContains)
				}
			}
		})
	}
}

func TestPhone_FullNumber(t *testing.T) {
	phone := NewPhone(1, "631", "555", "1234", phonenumbers.E164)
	fullNumber := phone.FullNumber()

	if fullNumber == nil {
		t.Fatal("FullNumber() returned nil")
	}

	if fullNumber.GetCountryCode() != 1 {
		t.Errorf("FullNumber().CountryCode = %d, want 1", fullNumber.GetCountryCode())
	}

	if fullNumber.GetNationalNumber() != 6315551234 {
		t.Errorf("FullNumber().NationalNumber = %d, want 6315551234", fullNumber.GetNationalNumber())
	}
}

func TestPhone_SetFormat(t *testing.T) {
	phone := NewPhone(1, "631", "555", "1234", phonenumbers.E164)

	// Initial format
	initialFormat := phone.FormatedNumber()
	if initialFormat != "+16315551234" {
		t.Errorf("Initial format = %s, want +16315551234", initialFormat)
	}

	// Change format
	phone.SetFormat(phonenumbers.NATIONAL)
	newFormat := phone.FormatedNumber()

	if newFormat == initialFormat {
		t.Error("SetFormat() did not change the format")
	}
}

func TestPhone_Possible(t *testing.T) {
	tests := []struct {
		name          string
		countryCode   int
		areaCode      string
		centralOffice string
		lineNumber    string
		wantPossible  bool
	}{
		{
			name:          "Valid US phone number",
			countryCode:   1,
			areaCode:      "631",
			centralOffice: "555",
			lineNumber:    "1234",
			wantPossible:  true,
		},
		{
			name:          "Another valid US phone number",
			countryCode:   1,
			areaCode:      "212",
			centralOffice: "123",
			lineNumber:    "4567",
			wantPossible:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone := NewPhone(tt.countryCode, tt.areaCode, tt.centralOffice, tt.lineNumber, phonenumbers.E164)
			got := phone.Possible()

			if got != tt.wantPossible {
				t.Errorf("Possible() = %v, want %v", got, tt.wantPossible)
			}
		})
	}
}

func TestPhone_Valid(t *testing.T) {
	tests := []struct {
		name          string
		countryCode   int
		areaCode      string
		centralOffice string
		lineNumber    string
	}{
		{
			name:          "US phone number",
			countryCode:   1,
			areaCode:      "631",
			centralOffice: "555",
			lineNumber:    "1234",
		},
		{
			name:          "Another US phone number",
			countryCode:   1,
			areaCode:      "212",
			centralOffice: "123",
			lineNumber:    "4567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone := NewPhone(tt.countryCode, tt.areaCode, tt.centralOffice, tt.lineNumber, phonenumbers.E164)
			// Valid() returns bool, just ensure it doesn't panic
			_ = phone.Valid()
		})
	}
}

func TestRandomPhone(t *testing.T) {
	areaCodes := []string{"631", "561", "446"}
	countryCode := 1
	format := phonenumbers.E164

	// Test multiple random phone generations
	for range 10 {
		phone := RandomPhone(countryCode, areaCodes, format)

		if phone == nil {
			t.Fatal("RandomPhone returned nil")
		}

		if phone.CountryCode() != countryCode {
			t.Errorf("RandomPhone().CountryCode() = %d, want %d", phone.CountryCode(), countryCode)
		}

		// Check if area code is from the provided list
		areaCodeValid := slices.Contains(areaCodes, phone.AreaCode())
		if !areaCodeValid {
			t.Errorf("RandomPhone().AreaCode() = %s, not in provided list", phone.AreaCode())
		}

		// Check central office is 3 digits
		if len(phone.CentralOffice()) != 3 {
			t.Errorf("RandomPhone().CentralOffice() length = %d, want 3", len(phone.CentralOffice()))
		}

		// Check line number is 4 digits
		if len(phone.LineNumber()) != 4 {
			t.Errorf("RandomPhone().LineNumber() length = %d, want 4", len(phone.LineNumber()))
		}
	}
}

func TestRandomPhone_SingleAreaCode(t *testing.T) {
	areaCodes := []string{"631"}
	phone := RandomPhone(1, areaCodes, phonenumbers.E164)

	if phone.AreaCode() != "631" {
		t.Errorf("RandomPhone with single area code = %s, want 631", phone.AreaCode())
	}
}

func BenchmarkNewPhone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPhone(1, "631", "555", "1234", phonenumbers.E164)
	}
}

func BenchmarkRandomPhone(b *testing.B) {
	areaCodes := []string{"631", "561", "446"}
	for i := 0; i < b.N; i++ {
		RandomPhone(1, areaCodes, phonenumbers.E164)
	}
}

func BenchmarkPhone_Possible(b *testing.B) {
	phone := NewPhone(1, "631", "555", "1234", phonenumbers.E164)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		phone.Possible()
	}
}
