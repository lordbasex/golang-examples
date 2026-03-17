//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	wasmhttp "github.com/nlepage/go-wasm-http-server"

	"certgen/internal/certgen"
	"certgen/internal/validation"
	"certgen/pkg/ziputil"
)

type certgenRequest struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
	Org    string `json:"org"`
}

type validationErrors map[string]string

// writeJSON writes a JSON response (used for errors and validation feedback).
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError is a small helper to keep error responses consistent.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]interface{}{
		"ok":      false,
		"message": msg,
	})
}

// writeValidation returns validation errors in a UI-friendly format.
func writeValidation(w http.ResponseWriter, errs validationErrors) {
	writeJSON(w, http.StatusBadRequest, map[string]interface{}{
		"ok":      false,
		"message": "validation_failed",
		"errors":  errs,
	})
}

// validateReq validates and normalizes input values.
// We do validation in the UI and again here (authoritative) to avoid generating broken certs.
func validateReq(req certgenRequest) validationErrors {
	req.Domain = strings.TrimSpace(req.Domain)
	req.IP = strings.TrimSpace(req.IP)
	req.Org = strings.TrimSpace(req.Org)

	errs := validation.ValidateAll(validation.Input{Domain: req.Domain, IP: req.IP, Org: req.Org})
	out := validationErrors(errs)
	return out
}

func certgenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[certgen] method_not_allowed method=%s path=%s", r.Method, r.URL.Path)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	defer r.Body.Close()

	var req certgenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[certgen] invalid_json err=%v", err)
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	log.Printf("[certgen] request domain=%q ip=%q org=%q", req.Domain, req.IP, req.Org)

	if errs := validateReq(req); len(errs) > 0 {
		log.Printf("[certgen] validation_failed errors=%v", errs)
		writeValidation(w, errs)
		return
	}

	out, err := certgen.Generate(certgen.Input{
		Domain: strings.TrimSpace(req.Domain),
		IP:     strings.TrimSpace(req.IP),
		Org:    strings.TrimSpace(req.Org),
	})
	if err != nil {
		log.Printf("[certgen] certgen_failed err=%v", err)
		writeError(w, http.StatusInternalServerError, "certgen_failed")
		return
	}

	zipBytes, err := ziputil.ZipFiles(out.Files)
	if err != nil {
		log.Printf("[certgen] zip_failed err=%v", err)
		writeError(w, http.StatusInternalServerError, "zip_failed")
		return
	}

	log.Printf("[certgen] ok zip_bytes=%d files=%d", len(zipBytes), len(out.Files))

	downloadName := buildZipFilename(strings.TrimSpace(req.Domain))
	log.Printf("[certgen] download_name=%q", downloadName)

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+downloadName+`"`)
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(zipBytes)
}

// buildZipFilename generates a stable filename:
// domain-YYYYMMDD-HHMMSS.zip
func buildZipFilename(domain string) string {
	now := time.Now()
	stamp := now.Format("20060102-150405")
	safeDomain := sanitizeFilenamePart(domain)
	if safeDomain == "" {
		safeDomain = "certs"
	}
	return safeDomain + "-" + stamp + ".zip"
}

// sanitizeFilenamePart replaces unsafe filename characters.
func sanitizeFilenamePart(value string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(strings.TrimSpace(value))
}

func main() {
	// All logs from WASM (stdout/stderr) show up in the browser console.
	log.Printf("Copyright © 2026 CNSoluciones - fpereira@cnsoluciones.com")
	log.Printf("[wasm] starting handlers under /api/* (Service Worker scope)")
	http.HandleFunc("/certgen", certgenHandler)
	wasmhttp.Serve(nil)
	select {}
}

