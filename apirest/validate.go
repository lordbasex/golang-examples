package main

import (
	"net/url"
	"strconv"
)

func (c *CreateCredit) validateCreateCreditCard() url.Values {
	errs := url.Values{}

	// Validate UniqueID (example: non-empty)
	if c.UniqueID == "" {
		errs.Add("uniqueid", "UniqueID is required")
	}

	// Validate Card (example: numeric and length)
	if c.Card == 0 {
		errs.Add("card", "Card number is required")
	}
	cardStr := strconv.FormatInt(c.Card, 10)
	if len(cardStr) < 13 || len(cardStr) > 19 {
		errs.Add("card", "Invalid card number")
	}

	// Validate ExpirationDate (example: numeric and length)
	if c.ExpirationDate == 0 {
		errs.Add("expiration_date", "Expiration date is required")
	}
	expirationStr := strconv.Itoa(c.ExpirationDate)
	if len(expirationStr) != 4 {
		errs.Add("expiration_date", "Invalid expiration date")
	}

	// Validate SecurityCode (example: numeric and up to 3 digits)
	if c.SecurityCode == 0 {
		errs.Add("security_code", "Security code is required")
	} else {
		securityCodeStr := strconv.Itoa(c.SecurityCode)
		if len(securityCodeStr) > 3 || !isNumeric(securityCodeStr) {
			errs.Add("security_code", "Invalid security code")
		}
	}

	// Validate Amount (example: non-negative)
	if c.Amount <= 0 {
		errs.Add("amount", "Amount must be greater than 0")
	}

	// Validate IDCustomer (example: non-negative)
	if c.IDCustomer <= 0 {
		errs.Add("id_customer", "Invalid IDCustomer")
	}

	return errs
}

func (s *Survey) validateCreateSurvey() url.Values {
	errs := url.Values{}

	// Validate UniqueID (example: non-empty)
	if s.UniqueID == "" {
		errs.Add("uniqueid", "UniqueID is required")
	}

	// Validate agent (example: numeric and length)
	if s.Agent <= 0 || s.Agent > 999999 {
		errs.Add("agent", "agent number is required and should be between 1 and 999999")
	}

	// Validate queue (example: numeric and length)
	if s.Queue <= 0 || s.Queue > 999999 {
		errs.Add("queue", "queue number is required and should be between 1 and 999999")
	}

	// Validate phone (example: numeric and length)
	if s.Phone <= 0 || s.Phone > 999999999999999999 {
		errs.Add("phone", "phone number is required and should be between 1 and 999999999999999999")
	}

	// Validate IDCustomer (example: numeric and length)
	if s.IDCustomer <= 0 || s.IDCustomer > 999999999999999999 {
		errs.Add("id_customer", "id_customer number is required and should be between 1 and 999999999999999999")
	}

	return errs
}
