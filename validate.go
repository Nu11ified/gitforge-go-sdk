package gitforge

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"
)

// ValidateWebhookOptions configures webhook validation.
type ValidateWebhookOptions struct {
	// Timestamp is the value of X-GitForge-Timestamp header (unix seconds).
	Timestamp string
	// Tolerance is the maximum age in seconds. nil = default 300s, Ptr(0) = skip freshness check.
	Tolerance *int
}

// ValidateWebhookSignature verifies the HMAC-SHA256 signature.
// The payload parameter should be the full signed content (which may include
// a "timestamp." prefix when replay protection is used).
func ValidateWebhookSignature(payload, signature, secret string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) == 1
}

// ValidateWebhook validates signature and optionally checks timestamp freshness.
//
// When opts.Timestamp is set, the signature is verified over "timestamp.payload"
// (Stripe-style replay protection). When empty, the signature is verified over
// the raw payload only (backward compatibility with old deliveries).
func ValidateWebhook(payload, secret, signature string, opts *ValidateWebhookOptions) bool {
	// Build the signed content the same way the server does.
	signedPayload := payload
	var timestamp string
	if opts != nil && opts.Timestamp != "" {
		timestamp = opts.Timestamp
		signedPayload = timestamp + "." + payload
	}

	if !ValidateWebhookSignature(signedPayload, signature, secret) {
		return false
	}

	if opts == nil {
		return true
	}

	tolerance := 300
	if opts.Tolerance != nil {
		tolerance = *opts.Tolerance
	}

	if timestamp != "" && tolerance > 0 {
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return false
		}
		now := time.Now().Unix()
		if math.Abs(float64(now-ts)) > float64(tolerance) {
			return false
		}
	}

	return true
}
