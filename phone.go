package main

import (
	"fmt"
	"math/rand"

	"github.com/nyaruka/phonenumbers"
)

type phone struct {
	countryCode   int
	areaCode      string
	centralOffice string
	lineNumber    string
	format        phonenumbers.PhoneNumberFormat
}

func NewPhone(countryCode int, areaCode, centralOffice, lineNumber string, format phonenumbers.PhoneNumberFormat) *phone {
	return &phone{countryCode, areaCode, centralOffice, lineNumber, format}
}

func RandomPhone(countryCode int, areaCodes []string, format phonenumbers.PhoneNumberFormat) *phone {
	return NewPhone(countryCode, areaCodes[rand.Intn(len(areaCodes))], fmt.Sprintf("%03d", rand.Intn(1000)), fmt.Sprintf("%04d", rand.Intn(10000)), format)
}

func (phone *phone) CountryCode() int {
	return phone.countryCode
}

func (phone *phone) AreaCode() string {
	return phone.areaCode
}

func (phone *phone) CentralOffice() string {
	return phone.centralOffice
}

func (phone *phone) LineNumber() string {
	return phone.lineNumber
}

func (phone *phone) FormatedNumber() string {
	return phonenumbers.Format(phone.FullNumber(), phone.format)
}

func (phone *phone) FullNumber() *phonenumbers.PhoneNumber {
	var concatenatedNumber = fmt.Sprintf("+%d%s%s%s", phone.countryCode, phone.areaCode, phone.centralOffice, phone.lineNumber)

	num, err := phonenumbers.Parse(concatenatedNumber, phonenumbers.GetRegionCodeForCountryCode(phone.countryCode))
	if err != nil {
		panic(err)
	}

	return num
}

func (phone *phone) SetFormat(format phonenumbers.PhoneNumberFormat) {
	phone.format = format
}

func (phone *phone) Possible() bool {
	return phonenumbers.IsPossibleNumber(phone.FullNumber())
}

func (phone *phone) Valid() bool {
	return phonenumbers.IsValidNumber(phone.FullNumber())
}
