package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordService_Hash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "securePassword123!",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "unicode password",
			password: "密码Пароль🔐",
			wantErr:  false,
		},
		{
			name:     "very long password",
			password: string(make([]byte, 72)), // bcrypt max is 72 bytes
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newPasswordService(bcrypt.MinCost) // Use MinCost for faster tests

			hash, err := svc.Hash(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if hash == "" {
					t.Error("Hash() returned empty hash")
				}
				if hash == tt.password {
					t.Error("Hash() returned plaintext password")
				}
			}
		})
	}
}

func TestPasswordService_Verify(t *testing.T) {
	svc := newPasswordService(bcrypt.MinCost)

	// Create a known hash
	password := "testPassword123!"
	hash, err := svc.Hash(password)
	if err != nil {
		t.Fatalf("Failed to create hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "wrongPassword",
			hash:     hash,
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			wantErr:  true,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			wantErr:  true,
		},
		{
			name:     "invalid hash",
			password: password,
			hash:     "not-a-valid-hash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Verify(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordService_CostBounds(t *testing.T) {
	tests := []struct {
		name         string
		cost         int
		expectedCost int
	}{
		{
			name:         "cost below minimum uses default",
			cost:         0,
			expectedCost: bcrypt.DefaultCost,
		},
		{
			name:         "cost at minimum",
			cost:         bcrypt.MinCost,
			expectedCost: bcrypt.MinCost,
		},
		{
			name:         "cost within range",
			cost:         12,
			expectedCost: 12,
		},
		{
			name:         "cost above maximum uses maximum",
			cost:         50,
			expectedCost: bcrypt.MaxCost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newPasswordService(tt.cost)
			if svc.cost != tt.expectedCost {
				t.Errorf("newPasswordService() cost = %v, want %v", svc.cost, tt.expectedCost)
			}
		})
	}
}

func TestPasswordService_HashUniqueness(t *testing.T) {
	svc := newPasswordService(bcrypt.MinCost)
	password := "samePassword123!"

	hash1, err := svc.Hash(password)
	if err != nil {
		t.Fatalf("First hash failed: %v", err)
	}

	hash2, err := svc.Hash(password)
	if err != nil {
		t.Fatalf("Second hash failed: %v", err)
	}

	// Each hash should be different (bcrypt uses random salt)
	if hash1 == hash2 {
		t.Error("Hash() should produce unique hashes for the same password")
	}

	// But both should verify correctly
	if err := svc.Verify(password, hash1); err != nil {
		t.Errorf("First hash verification failed: %v", err)
	}
	if err := svc.Verify(password, hash2); err != nil {
		t.Errorf("Second hash verification failed: %v", err)
	}
}

func BenchmarkPasswordService_Hash(b *testing.B) {
	svc := newPasswordService(12) // Production cost

	for b.Loop() {
		_, _ = svc.Hash("benchmarkPassword123!")
	}
}

func BenchmarkPasswordService_Verify(b *testing.B) {
	svc := newPasswordService(12)
	hash, _ := svc.Hash("benchmarkPassword123!")

	b.ResetTimer()
	for b.Loop() {
		_ = svc.Verify("benchmarkPassword123!", hash)
	}
}
