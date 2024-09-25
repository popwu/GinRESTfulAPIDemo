package utils

import (
	"os"
	"testing"
)

// 辅助函数：设置环境变量
func setEnv(key, value string) func() {
	oldValue := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		os.Setenv(key, oldValue)
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		nickname    string
		signingKey  string
		expectError bool
	}{
		{"Valid token", 1, "user1", "test_secret", false},
		{"Missing signing key", 1, "user1", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setEnv("JWT_SIGNING_KEY", tt.signingKey)
			defer cleanup()

			token, err := GenerateToken(tt.userID, tt.nickname)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token == "" {
					t.Errorf("Expected a token, but got an empty string")
				}
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		nickname    string
		signingKey  string
		parseKey    string
		expectError bool
	}{
		{"Valid token", 1, "user1", "test_secret", "test_secret", false},
		{"Invalid signing key", 1, "user1", "test_secret", "wrong_secret", true},
		{"Missing signing key", 1, "user1", "test_secret", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 生成 token
			cleanup := setEnv("JWT_SIGNING_KEY", tt.signingKey)
			token, err := GenerateToken(tt.userID, tt.nickname)
			if err != nil {
				t.Fatalf("Failed to generate token: %v", err)
			}
			cleanup()

			// 解析 token
			cleanup = setEnv("JWT_SIGNING_KEY", tt.parseKey)
			defer cleanup()

			userID, nickname, err := ParseToken(token)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if userID != tt.userID {
					t.Errorf("Expected userID %d, but got %d", tt.userID, userID)
				}
				if nickname != tt.nickname {
					t.Errorf("Expected nickname %s, but got %s", tt.nickname, nickname)
				}
			}
		})
	}
}
