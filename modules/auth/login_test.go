package auth

import (
	"github.com/tryffel/market/modules/Error"
	"testing"
	"time"
)

func TestPlainToken(t *testing.T) {
	privateKey := "abcdefgh123"
	userId := "user1"
	nonce := "very-secret-nonce"
	expire := time.Second * 30

	token, err := NewToken(userId, privateKey, nonce, expire)
	if err != nil {
		t.Error("Failed to issue new token: ", err)
	}

	user, non, err := ValidateToken(token, privateKey, true, time.Second*60)
	if user != userId {
		t.Errorf("Invalid user_id, expected %s, got %s", userId, user)
	}
	if non != nonce {
		t.Errorf("Invalid nonce, expected %s, got %s", nonce, non)
	}
	if err != nil {
		t.Error("Failed to validate token: ", err)
	}
}

func TestExpiredToken(t *testing.T) {
	privateKey := "abcdefgh123"
	userId := "user1"
	nonce := "very-secret-nonce"
	expire := time.Second * 2

	token, err := NewToken(userId, privateKey, nonce, expire)
	if err != nil {
		t.Error("Failed to issue new token: ", err)
	}

	time.Sleep(time.Second * 1)
	user, _, err := ValidateToken(token, privateKey, true, time.Millisecond*1)
	if err != nil {
		t.Error("error", err)
	}
	if user != userId {
		t.Errorf("Invalid user_id, expected %s, got %s", userId, user)
	}

	time.Sleep(time.Second * 3)
	user, _, err = ValidateToken(token, privateKey, true, time.Millisecond*1)
	if err == nil {
		t.Error("No error!")
	}
	if user != "" {
		t.Errorf("Invalid user_id, expected %s, got %s", "", user)
	}

	errorMsg := err.(*Error.Error).EndUserMessage()
	if errorMsg != "token expired" {
		t.Error("Got invalid error msg: ", errorMsg)
	}
}

func TestInvalidToken(t *testing.T) {
	privateKey := "abcdefgh123"
	userId := "user1"
	nonce := "very-secret-nonce"
	expire := time.Second * 2

	token, err := NewToken(userId, privateKey, nonce, expire)
	if err != nil {
		t.Error("Failed to issue new token: ", err)
	}

	newPrivateKey := "abcde"

	user, _, err := ValidateToken(token, newPrivateKey, true, time.Millisecond*1)
	if err == nil {
		t.Error("No error, expected 'invalid token'")
	}
	if user != "" {
		t.Errorf("Invalid user_id, expected %s, got %s", "", user)
	}
}
