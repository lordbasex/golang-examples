package validation

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
)

// Errors is a field -> error message map.
// This lets the backend report multiple validation problems at once.
type Errors map[string]string

type Input struct {
	Domain string
	IP     string
	Org    string
}

// Validate validates all fields and returns a single error describing the first failure.
// We validate in both the UI and the WASM backend so the user gets fast feedback,
// but the backend stays authoritative.
func Validate(in Input) error {
	// Keep this function for callers that only want one error.
	// For full field-level validation, prefer ValidateAll.
	errs := ValidateAll(in)
	if len(errs) == 0 {
		return nil
	}

	// Return a stable "first" error to avoid randomness from map iteration.
	keys := make([]string, 0, len(errs))
	for k := range errs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return errors.New(errs[keys[0]])
}

// ValidateAll validates all fields and returns errors for every invalid field.
func ValidateAll(in Input) Errors {
	errs := Errors{}

	domain := strings.TrimSpace(in.Domain)
	ip := strings.TrimSpace(in.IP)
	org := strings.TrimSpace(in.Org)

	// Domain
	if domain == "" {
		errs["domain"] = "domain is required"
	} else if !isValidDomain(domain) {
		errs["domain"] = fmt.Sprintf("invalid domain: %q", domain)
	}

	// IP
	if ip == "" {
		errs["ip"] = "ip is required"
	} else if net.ParseIP(ip) == nil {
		errs["ip"] = fmt.Sprintf("invalid ip: %q", ip)
	}

	// Org
	if org == "" {
		errs["org"] = "org is required"
	} else if len(org) > 128 {
		errs["org"] = "org is too long (max 128 chars)"
	}

	return errs
}

func isValidDomain(domain string) bool {
	// Basic rules:
	// - at least 2 labels (pbx.local)
	// - labels: 1..63, only [a-zA-Z0-9-], cannot start/end with '-'
	// - total length <= 253
	if len(domain) > 253 {
		return false
	}
	if strings.Contains(domain, " ") {
		return false
	}
	domain = strings.Trim(domain, ".")
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return false
	}
	for _, label := range parts {
		if label == "" || len(label) > 63 {
			return false
		}
		for i := 0; i < len(label); i++ {
			c := label[i]
			isAlphaNum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
			if isAlphaNum || c == '-' {
				continue
			}
			return false
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
	}
	return true
}

