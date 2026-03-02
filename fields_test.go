package main

import (
	"strings"
	"testing"
)

func TestAddressCache_GetAddress(t *testing.T) {
	cache := &addressCache{}

	// First call should create a new address
	addr1 := cache.getAddress()
	if addr1 == nil {
		t.Fatal("getAddress() returned nil")
	}

	// Second call should return the same cached address
	addr2 := cache.getAddress()
	if addr2 == nil {
		t.Fatal("getAddress() returned nil on second call")
	}

	if addr1 != addr2 {
		t.Error("getAddress() did not return cached address")
	}

	// Verify address has expected fields populated
	if addr1.Address == "" {
		t.Error("Address field is empty")
	}
	if addr1.City == "" {
		t.Error("City field is empty")
	}
	if addr1.State == "" {
		t.Error("State field is empty")
	}
	if addr1.PostalCode == "" {
		t.Error("PostalCode field is empty")
	}
}

func TestAddressCache_Reset(t *testing.T) {
	cache := &addressCache{}

	// Get an address
	addr1 := cache.getAddress()
	if addr1 == nil {
		t.Fatal("getAddress() returned nil")
	}

	// Reset the cache
	cache.reset()

	if cache.address != nil {
		t.Error("reset() did not clear the cache")
	}

	// Get a new address after reset
	addr2 := cache.getAddress()
	if addr2 == nil {
		t.Fatal("getAddress() returned nil after reset")
	}

	// Addresses should be different after reset
	if addr1 == addr2 {
		t.Error("reset() did not clear the cached address reference")
	}
}

func TestGetFieldRegistry(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	expectedFields := []string{
		"firstname", "lastname", "address", "city", "state", "zip",
		"latitude", "longitude", "email", "dob", "username", "company",
	}

	// Check that all expected fields exist
	for _, fieldName := range expectedFields {
		field, exists := registry[fieldName]
		if !exists {
			t.Errorf("Expected field %s not found in registry", fieldName)
			continue
		}

		if field.Name == "" {
			t.Errorf("Field %s has empty Name", fieldName)
		}

		if field.Generator == nil {
			t.Errorf("Field %s has nil Generator", fieldName)
		}
	}

	// Verify the registry has the correct number of fields
	if len(registry) != len(expectedFields) {
		t.Errorf("Registry has %d fields, expected %d", len(registry), len(expectedFields))
	}
}

func TestGetFieldRegistry_Generators(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	tests := []struct {
		name      string
		fieldKey  string
		checkFunc func(string) bool
	}{
		{
			name:     "firstname generates non-empty string",
			fieldKey: "firstname",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "lastname generates non-empty string",
			fieldKey: "lastname",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "email generates non-empty string",
			fieldKey: "email",
			checkFunc: func(s string) bool {
				return len(s) > 0 && strings.Contains(s, "@")
			},
		},
		{
			name:     "address generates non-empty string",
			fieldKey: "address",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "city generates non-empty string",
			fieldKey: "city",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "state generates non-empty string",
			fieldKey: "state",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "zip generates non-empty string",
			fieldKey: "zip",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "username generates non-empty string",
			fieldKey: "username",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "company generates non-empty string",
			fieldKey: "company",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "dob generates non-empty string",
			fieldKey: "dob",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "latitude generates numeric string",
			fieldKey: "latitude",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
		{
			name:     "longitude generates numeric string",
			fieldKey: "longitude",
			checkFunc: func(s string) bool {
				return len(s) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, exists := registry[tt.fieldKey]
			if !exists {
				t.Fatalf("Field %s not found in registry", tt.fieldKey)
			}

			result := field.Generator()
			if !tt.checkFunc(result) {
				t.Errorf("Generator for %s produced invalid result: %s", tt.fieldKey, result)
			}
		})
	}
}

func TestGetFieldRegistry_AddressCaching(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	// Generate multiple address-related fields
	addr1 := registry["address"].Generator()
	city1 := registry["city"].Generator()

	// They should use the same cached address
	addr2 := registry["address"].Generator()
	city2 := registry["city"].Generator()

	if addr1 != addr2 {
		t.Error("Address should be cached between calls")
	}
	if city1 != city2 {
		t.Error("City should be cached between calls")
	}

	// After reset, should get different values
	cache.reset()
	addr3 := registry["address"].Generator()

	if addr1 == addr3 {
		t.Error("Address should be different after cache reset")
	}
}

func TestParseRequestedFields_EmptyString(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	fields, err := parseRequestedFields("", registry)

	if err != nil {
		t.Errorf("parseRequestedFields(\"\") returned error: %v", err)
	}

	if fields != nil {
		t.Errorf("parseRequestedFields(\"\") = %v, want nil", fields)
	}
}

func TestParseRequestedFields_SingleField(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	fields, err := parseRequestedFields("firstname", registry)

	if err != nil {
		t.Fatalf("parseRequestedFields(\"firstname\") returned error: %v", err)
	}

	if len(fields) != 1 {
		t.Fatalf("parseRequestedFields(\"firstname\") returned %d fields, want 1", len(fields))
	}

	if fields[0].Name != "First Name" {
		t.Errorf("Field name = %s, want First Name", fields[0].Name)
	}
}

func TestParseRequestedFields_MultipleFields(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	fields, err := parseRequestedFields("firstname,lastname,email", registry)

	if err != nil {
		t.Fatalf("parseRequestedFields() returned error: %v", err)
	}

	if len(fields) != 3 {
		t.Fatalf("parseRequestedFields() returned %d fields, want 3", len(fields))
	}

	expectedNames := []string{"First Name", "Last Name", "Email"}
	for i, expected := range expectedNames {
		if fields[i].Name != expected {
			t.Errorf("Field[%d].Name = %s, want %s", i, fields[i].Name, expected)
		}
	}
}

func TestParseRequestedFields_CaseInsensitive(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	tests := []string{
		"FIRSTNAME",
		"FirstName",
		"firstName",
		"FiRsTnAmE",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			fields, err := parseRequestedFields(input, registry)

			if err != nil {
				t.Fatalf("parseRequestedFields(%s) returned error: %v", input, err)
			}

			if len(fields) != 1 {
				t.Fatalf("parseRequestedFields(%s) returned %d fields, want 1", input, len(fields))
			}

			if fields[0].Name != "First Name" {
				t.Errorf("Field name = %s, want First Name", fields[0].Name)
			}
		})
	}
}

func TestParseRequestedFields_WithWhitespace(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	fields, err := parseRequestedFields("  firstname  ,  lastname  ,  email  ", registry)

	if err != nil {
		t.Fatalf("parseRequestedFields() with whitespace returned error: %v", err)
	}

	if len(fields) != 3 {
		t.Fatalf("parseRequestedFields() returned %d fields, want 3", len(fields))
	}
}

func TestParseRequestedFields_EmptyValues(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	fields, err := parseRequestedFields("firstname,,lastname", registry)

	if err != nil {
		t.Fatalf("parseRequestedFields() with empty value returned error: %v", err)
	}

	// Should skip empty values and return 2 fields
	if len(fields) != 2 {
		t.Fatalf("parseRequestedFields() returned %d fields, want 2", len(fields))
	}
}

func TestParseRequestedFields_UnknownField(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	_, err := parseRequestedFields("unknown_field", registry)

	if err == nil {
		t.Fatal("parseRequestedFields() with unknown field should return error")
	}

	expectedError := "unknown field: unknown_field"
	if err.Error() != expectedError {
		t.Errorf("Error message = %s, want %s", err.Error(), expectedError)
	}
}

func TestParseRequestedFields_MixedValidInvalid(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	_, err := parseRequestedFields("firstname,invalid_field,lastname", registry)

	if err == nil {
		t.Fatal("parseRequestedFields() with invalid field should return error")
	}

	if !strings.Contains(err.Error(), "unknown field") {
		t.Errorf("Error should mention unknown field, got: %s", err.Error())
	}
}

func TestParseRequestedFields_AllFields(t *testing.T) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)

	allFields := "firstname,lastname,address,city,state,zip,latitude,longitude,email,dob,username,company"
	fields, err := parseRequestedFields(allFields, registry)

	if err != nil {
		t.Fatalf("parseRequestedFields() with all fields returned error: %v", err)
	}

	if len(fields) != 12 {
		t.Fatalf("parseRequestedFields() returned %d fields, want 12", len(fields))
	}

	// Verify all generators work
	for i, field := range fields {
		result := field.Generator()
		if result == "" {
			t.Errorf("Field[%d] (%s) generator returned empty string", i, field.Name)
		}
	}
}

func BenchmarkGetFieldRegistry(b *testing.B) {
	cache := &addressCache{}
	for i := 0; i < b.N; i++ {
		getFieldRegistry(cache)
	}
}

func BenchmarkParseRequestedFields(b *testing.B) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)
	for i := 0; i < b.N; i++ {
		_, _ = parseRequestedFields("firstname,lastname,email", registry)
	}
}

func BenchmarkFieldGeneration(b *testing.B) {
	cache := &addressCache{}
	registry := getFieldRegistry(cache)
	field := registry["firstname"]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		field.Generator()
	}
}
