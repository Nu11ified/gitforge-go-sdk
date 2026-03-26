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
	Timestamp     string
	MaxAgeSeconds *int // nil = default 300s, Ptr(0) = skip timestamp check, Ptr(N) = N seconds
}

// ValidateWebhookSignature verifies the HMAC-SHA256 signature.
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
func ValidateWebhook(payload, secret, signature string, opts *ValidateWebhookOptions) bool {
	if !ValidateWebhookSignature(payload, signature, secret) {
		return false
	}

	if opts == nil {
		return true
	}

	maxAge := 300
	if opts.MaxAgeSeconds != nil {
		maxAge = *opts.MaxAgeSeconds
	}

	if opts.Timestamp != "" && maxAge > 0 {
		ts, err := strconv.ParseInt(opts.Timestamp, 10, 64)
		if err != nil {
			return false
		}
		now := time.Now().Unix()
		if math.Abs(float64(now-ts)) > float64(maxAge) {
			return false
		}
	}

	return true
}
