package gitforge

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func sign(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func intPtr(v int) *int {
	return &v
}

func TestValidateWebhookSignature_Valid(t *testing.T) {
	payload := `{"event":"push"}`
	sig := sign(payload, "my-secret")
	if !ValidateWebhookSignature(payload, sig, "my-secret") {
		t.Fatal("expected valid signature")
	}
}

func TestValidateWebhookSignature_TamperedPayload(t *testing.T) {
	sig := sign(`{"event":"push"}`, "my-secret")
	if ValidateWebhookSignature(`{"event":"hack"}`, sig, "my-secret") {
		t.Fatal("expected invalid for tampered payload")
	}
}

func TestValidateWebhookSignature_WrongSecret(t *testing.T) {
	sig := sign(`{"event":"push"}`, "correct")
	if ValidateWebhookSignature(`{"event":"push"}`, sig, "wrong") {
		t.Fatal("expected invalid for wrong secret")
	}
}

func TestValidateWebhookSignature_MissingPrefix(t *testing.T) {
	mac := hmac.New(sha256.New, []byte("my-secret"))
	mac.Write([]byte(`{"event":"push"}`))
	rawHex := hex.EncodeToString(mac.Sum(nil))
	if ValidateWebhookSignature(`{"event":"push"}`, rawHex, "my-secret") {
		t.Fatal("expected invalid without sha256= prefix")
	}
}

func TestValidateWebhook_WithTimestamp(t *testing.T) {
	payload := `{"event":"push"}`
	sig := sign(payload, "my-secret")
	ts := fmt.Sprintf("%d", time.Now().Unix())
	if !ValidateWebhook(payload, "my-secret", sig, &ValidateWebhookOptions{Timestamp: ts}) {
		t.Fatal("expected valid with fresh timestamp (default 300s)")
	}
}

func TestValidateWebhook_ExpiredTimestamp(t *testing.T) {
	payload := `{"event":"push"}`
	sig := sign(payload, "my-secret")
	ts := fmt.Sprintf("%d", time.Now().Unix()-600)
	if ValidateWebhook(payload, "my-secret", sig, &ValidateWebhookOptions{Timestamp: ts, MaxAgeSeconds: intPtr(300)}) {
		t.Fatal("expected invalid for expired timestamp")
	}
}

func TestValidateWebhook_SkipTimestampWhenZero(t *testing.T) {
	payload := `{"event":"push"}`
	sig := sign(payload, "my-secret")
	ts := fmt.Sprintf("%d", time.Now().Unix()-99999)
	if !ValidateWebhook(payload, "my-secret", sig, &ValidateWebhookOptions{Timestamp: ts, MaxAgeSeconds: intPtr(0)}) {
		t.Fatal("expected valid when maxAge is Ptr(0)")
	}
}

func TestValidateWebhook_NoTimestamp(t *testing.T) {
	payload := `{"event":"push"}`
	sig := sign(payload, "my-secret")
	if !ValidateWebhook(payload, "my-secret", sig, nil) {
		t.Fatal("expected valid without timestamp")
	}
}
