//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/nlepage/go-wasm-http-server"
)

type User struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Document  string `json:"document"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zipCode"`
	Country   string `json:"country"`
	Company   string `json:"company"`
	Role      string `json:"role"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type ValidationErrors map[string]string

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeValidation(w http.ResponseWriter, errors ValidationErrors) {
	writeJSON(w, http.StatusBadRequest, map[string]interface{}{
		"ok":     false,
		"message": "validation_failed",
		"errors": errors,
	})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]interface{}{
		"ok": false,
		"message": msg,
	})
}

func jsValueToString(v js.Value) string {
	if v.Type() == js.TypeString {
		return v.String()
	}
	return js.Global().Get("JSON").Call("stringify", v).String()
}

func awaitPromise(promise js.Value) (js.Value, error) {
	okCh := make(chan js.Value, 1)
	errCh := make(chan js.Value, 1)

	then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			okCh <- args[0]
		} else {
			okCh <- js.Undefined()
		}
		return nil
	})
	catcher := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			errCh <- args[0]
		} else {
			errCh <- js.ValueOf("promise_rejected")
		}
		return nil
	})
	promise.Call("then", then).Call("catch", catcher)

	var v js.Value
	var err error
	select {
	case v = <-okCh:
	case e := <-errCh:
		err = fmt.Errorf("%s", jsValueToString(e))
	}

	then.Release()
	catcher.Release()
	return v, err
}

func writePromiseJSON(w http.ResponseWriter, promise js.Value, successStatus int) {
	payload, err := awaitPromise(promise)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error")
		return
	}

	body := []byte(jsValueToString(payload))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(successStatus)
	_, _ = w.Write(body)
}

func validateUser(u User) ValidationErrors {
	errors := ValidationErrors{}

	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
	u.Phone = strings.TrimSpace(u.Phone)
	u.Document = strings.TrimSpace(u.Document)
	u.Address = strings.TrimSpace(u.Address)
	u.City = strings.TrimSpace(u.City)
	u.State = strings.TrimSpace(u.State)
	u.ZipCode = strings.TrimSpace(u.ZipCode)
	u.Country = strings.TrimSpace(u.Country)
	u.Company = strings.TrimSpace(u.Company)
	u.Role = strings.TrimSpace(u.Role)
	u.Notes = strings.TrimSpace(u.Notes)

	if u.FirstName == "" {
		errors["firstName"] = "First name is required"
	} else if len([]rune(u.FirstName)) > 80 {
		errors["firstName"] = "Maximum 80 characters"
	}

	if u.LastName == "" {
		errors["lastName"] = "Last name is required"
	} else if len([]rune(u.LastName)) > 80 {
		errors["lastName"] = "Maximum 80 characters"
	}

	if u.Email == "" {
		errors["email"] = "Email is required"
	} else if !isValidEmail(u.Email) {
		errors["email"] = "Invalid email"
	} else if len([]rune(u.Email)) > 120 {
		errors["email"] = "Maximum 120 characters"
	}

	if u.Phone != "" {
		if !isValidPhone(u.Phone) {
			errors["phone"] = "Invalid phone"
		} else if len([]rune(u.Phone)) > 30 {
			errors["phone"] = "Maximum 30 characters"
		}
	}

	if len([]rune(u.Document)) > 30 {
		errors["document"] = "Maximum 30 characters"
	}
	if len([]rune(u.Address)) > 180 {
		errors["address"] = "Maximum 180 characters"
	}
	if len([]rune(u.City)) > 80 {
		errors["city"] = "Maximum 80 characters"
	}
	if len([]rune(u.State)) > 80 {
		errors["state"] = "Maximum 80 characters"
	}
	if len([]rune(u.ZipCode)) > 20 {
		errors["zipCode"] = "Maximum 20 characters"
	}
	if len([]rune(u.Country)) > 80 {
		errors["country"] = "Maximum 80 characters"
	}
	if len([]rune(u.Company)) > 120 {
		errors["company"] = "Maximum 120 characters"
	}
	if len([]rune(u.Role)) > 80 {
		errors["role"] = "Maximum 80 characters"
	}
	if len([]rune(u.Notes)) > 1000 {
		errors["notes"] = "Maximum 1000 characters"
	}

	return errors
}

func isValidEmail(value string) bool {
	if strings.Contains(value, " ") {
		return false
	}
	parts := strings.Split(value, "@")
	if len(parts) != 2 {
		return false
	}
	local := strings.TrimSpace(parts[0])
	domain := strings.TrimSpace(parts[1])
	if local == "" || domain == "" {
		return false
	}
	if !strings.Contains(domain, ".") {
		return false
	}
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}
	return true
}

func isValidPhone(value string) bool {
	countDigits := 0
	for _, r := range value {
		if r >= '0' && r <= '9' {
			countDigits++
			continue
		}
		switch r {
		case '+', '-', ' ', '(', ')':
		default:
			return false
		}
	}
	return countDigits >= 6 && countDigits <= 20
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		promise := js.Global().Call("apiListUsers")
		writePromiseJSON(w, promise, http.StatusOK)

	case http.MethodPost:
		defer r.Body.Close()

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json")
			return
		}

		if errors := validateUser(user); len(errors) > 0 {
			writeValidation(w, errors)
			return
		}

		now := time.Now().UTC().Format(time.RFC3339)
		user.CreatedAt = now
		user.UpdatedAt = now

		jsObj := js.ValueOf(map[string]interface{}{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"phone":     user.Phone,
			"document":  user.Document,
			"address":   user.Address,
			"city":      user.City,
			"state":     user.State,
			"zipCode":   user.ZipCode,
			"country":   user.Country,
			"company":   user.Company,
			"role":      user.Role,
			"notes":     user.Notes,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		})

		promise := js.Global().Call("apiCreateUser", jsObj)
		writePromiseJSON(w, promise, http.StatusCreated)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func userDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	idStr = strings.Trim(idStr, "/")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "missing_user_id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_user_id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		promise := js.Global().Call("apiGetUser", id)
		writePromiseJSON(w, promise, http.StatusOK)

	case http.MethodPut:
		defer r.Body.Close()

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_json")
			return
		}

		user.ID = id

		if errors := validateUser(user); len(errors) > 0 {
			writeValidation(w, errors)
			return
		}

		if user.CreatedAt == "" {
			user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		}
		user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

		jsObj := js.ValueOf(map[string]interface{}{
			"id":        user.ID,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"phone":     user.Phone,
			"document":  user.Document,
			"address":   user.Address,
			"city":      user.City,
			"state":     user.State,
			"zipCode":   user.ZipCode,
			"country":   user.Country,
			"company":   user.Company,
			"role":      user.Role,
			"notes":     user.Notes,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		})

		promise := js.Global().Call("apiUpdateUser", jsObj)
		writePromiseJSON(w, promise, http.StatusOK)

	case http.MethodDelete:
		promise := js.Global().Call("apiDeleteUser", id)
		writePromiseJSON(w, promise, http.StatusOK)

	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userDetailHandler)
	wasmhttp.Serve(nil)
	select {}
}
