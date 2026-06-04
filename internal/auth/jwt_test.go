package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJwt(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	tokenSecret1 := "test-secret-1"
	tokenSecret2 := "test-secret-2"
	expiresIn := time.Duration(24 * time.Hour)
	expired := time.Duration(-1 * time.Second)
	jwt1, _ := MakeJWT(userID1, tokenSecret1, expiresIn)
	jwt2, _ := MakeJWT(userID2, tokenSecret2, expiresIn)
	expiredToken, _ := MakeJWT(userID2, tokenSecret2, expired)

	tests := []struct {
		name        string
		userid      uuid.UUID
		jwt         string
		tokenSecret string
		wantErr     bool
	}{
		{
			name:        "Correct jwt",
			userid:      userID1,
			jwt:         jwt1,
			tokenSecret: tokenSecret1,
			wantErr:     false,
		},
		{
			name:        "jwt is not valid",
			userid:      userID1,
			jwt:         jwt2,
			tokenSecret: tokenSecret1,
			wantErr:     true,
		},
		{
			name:        "Empty jwt",
			userid:      uuid.Nil,
			jwt:         "not.a.jwt",
			tokenSecret: tokenSecret1,
			wantErr:     true,
		},
		{
			name:        "Invalid jwt",
			userid:      userID1,
			jwt:         jwt1,
			tokenSecret: "invalidtokenSecret",
			wantErr:     true,
		},
		{
			name:        "Expired jwt",
			userid:      userID2,
			jwt:         expiredToken,
			tokenSecret: tokenSecret2,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := ValidateJWT(tt.jwt, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && valid != tt.userid {
				t.Errorf("ValidateJWT() expects %v, got %v", tt.userid, valid)
			}
		})
	}
}
