package main

import (
	"fmt"
	"strings"

	"github.com/go-faker/faker/v4"
)

type Field struct {
	Name      string
	Generator func() string
}

type addressCache struct {
	address *faker.RealAddress
}

func (a *addressCache) getAddress() *faker.RealAddress {
	if a.address == nil {
		var addr = faker.GetRealAddress()
		a.address = &addr
	}
	return a.address
}

func (a *addressCache) reset() {
	a.address = nil
}

func getFieldRegistry(addrCache *addressCache) map[string]Field {
	return map[string]Field{
		"firstname": {Name: "First Name", Generator: func() string { return faker.FirstName() }},
		"lastname":  {Name: "Last Name", Generator: func() string { return faker.LastName() }},
		"address":   {Name: "Address", Generator: func() string { return addrCache.getAddress().Address }},
		"city":      {Name: "City", Generator: func() string { return addrCache.getAddress().City }},
		"state":     {Name: "State", Generator: func() string { return addrCache.getAddress().State }},
		"zip":       {Name: "Zip", Generator: func() string { return addrCache.getAddress().PostalCode }},
		"latitude":  {Name: "Latitude", Generator: func() string { return fmt.Sprintf("%f", addrCache.getAddress().Coordinates.Latitude) }},
		"longitude": {Name: "Longitude", Generator: func() string { return fmt.Sprintf("%f", addrCache.getAddress().Coordinates.Longitude) }},
		"email":     {Name: "Email", Generator: func() string { email := faker.Email(); return email }},
		"dob":       {Name: "Date of Birth", Generator: func() string { date := faker.Date(); return date }},
		"username":  {Name: "Username", Generator: func() string { username := faker.Username(); return username }},
		"company":   {Name: "Company", Generator: func() string { domain := faker.DomainName(); return domain }},
		// Add more fields here as needed
	}
}

func parseRequestedFields(fieldNames string, registry map[string]Field) ([]Field, error) {
	if fieldNames == "" {
		return nil, nil
	}

	var fields []Field
	requestedNames := strings.SplitSeq(fieldNames, ",")

	for name := range requestedNames {
		name = strings.TrimSpace(strings.ToLower(name))
		if name == "" {
			continue
		}

		field, exists := registry[name]
		if !exists {
			return nil, fmt.Errorf("unknown field: %s", name)
		}
		fields = append(fields, field)
	}

	return fields, nil
}
