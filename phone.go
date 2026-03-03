package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/nyaruka/phonenumbers"
)

type Phone struct {
	countryCode   int
	areaCode      string
	centralOffice string
	lineNumber    string
	format        phonenumbers.PhoneNumberFormat
}

func NewPhone(countryCode int, areaCode, centralOffice, lineNumber string, format phonenumbers.PhoneNumberFormat) (*Phone, error) {
	phone := &Phone{countryCode, areaCode, centralOffice, lineNumber, format}

	if !phone.ValidCentralOffice() {
		return nil, fmt.Errorf("invalid central office code: %s", centralOffice)
	}

	return phone, nil
}

func RandomPhone(countryCode int, areaCodes []string, format phonenumbers.PhoneNumberFormat) *Phone {
	// Generate valid NXX code: first digit 2-9, avoid N11 patterns
	var centralOffice string
	for {
		n := rand.Intn(8) + 2 // 2-9
		x1 := rand.Intn(10)   // 0-9
		x2 := rand.Intn(10)   // 0-9

		// Avoid N11 patterns
		if x1 == 1 && x2 == 1 {
			continue
		}

		centralOffice = fmt.Sprintf("%d%d%d", n, x1, x2)
		break
	}

	phone, err := NewPhone(countryCode, areaCodes[rand.Intn(len(areaCodes))], centralOffice, fmt.Sprintf("%04d", rand.Intn(10000)), format)
	if err != nil {
		panic(err)
	}
	return phone
}

func (phone *Phone) CountryCode() int {
	return phone.countryCode
}

func (phone *Phone) AreaCode() string {
	return phone.areaCode
}

func (phone *Phone) CentralOffice() string {
	return phone.centralOffice
}

func (phone *Phone) LineNumber() string {
	return phone.lineNumber
}

func (phone *Phone) FormatedNumber() string {
	return phonenumbers.Format(phone.FullNumber(), phone.format)
}

func (phone *Phone) FullNumber() *phonenumbers.PhoneNumber {
	var concatenatedNumber = fmt.Sprintf("+%d%s%s%s", phone.countryCode, phone.areaCode, phone.centralOffice, phone.lineNumber)

	num, err := phonenumbers.Parse(concatenatedNumber, phonenumbers.GetRegionCodeForCountryCode(phone.countryCode))
	if err != nil {
		panic(err)
	}

	return num
}

func (phone *Phone) SetFormat(format phonenumbers.PhoneNumberFormat) {
	phone.format = format
}

func (phone *Phone) Possible() bool {
	return phonenumbers.IsPossibleNumber(phone.FullNumber())
}

func (phone *Phone) Valid() bool {
	return phonenumbers.IsValidNumber(phone.FullNumber())
}

func (phone *Phone) ValidCentralOffice() bool {
	// NXX code must be exactly 3 digits
	if len(phone.centralOffice) != 3 {
		return false
	}

	// Parse each digit
	n, err1 := strconv.Atoi(string(phone.centralOffice[0]))
	x1, err2 := strconv.Atoi(string(phone.centralOffice[1]))
	x2, err3 := strconv.Atoi(string(phone.centralOffice[2]))

	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	// First digit (N) must be 2-9
	if n < 2 || n > 9 {
		return false
	}

	// Check for invalid N11 patterns (211, 311, 411, 511, 611, 711, 811, 911)
	if x1 == 1 && x2 == 1 {
		return false
	}

	return true
}
