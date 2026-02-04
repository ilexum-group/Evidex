// Package utils provides utility functions for Evidex
package utils

import (
	"crypto/md5"  // #nosec G501 - MD5 used for forensic verification, not security
	"crypto/sha1" // #nosec G505 - SHA1 used for forensic verification, not security
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

// CalculateSHA1 calculates SHA-1 hash of a file.
func CalculateSHA1(filePath string) (string, error) {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	hash := sha1.New() // #nosec G401 - SHA1 used for forensic compatibility, not cryptographic security
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CalculateSHA256 calculates SHA-256 hash of a file.
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CalculateSHA512 calculates SHA-512 hash of a file.
func CalculateSHA512(filePath string) (string, error) {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CalculateMD5 calculates MD5 hash of a file (for compatibility only).
func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	hash := md5.New() // #nosec G401 - MD5 used for forensic compatibility, not cryptographic security
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// IsReadOnly checks if a file can be opened in read-only mode.
func IsReadOnly(filePath string) (bool, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0) // #nosec G304 - filepath is user-provided evidence path
	if err != nil {
		if os.IsPermission(err) {
			return false, fmt.Errorf("permission denied")
		}
		return false, err
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()
	return true, nil
}

// GenerateEvidenceID creates a unique evidence package identifier.
func GenerateEvidenceID(hostname string) string {
	randomID := GenerateRandomID()
	return fmt.Sprintf("BITX-%s-%s", hostname, randomID)
}

// GenerateRandomID creates a random identifier string UUID-like.
func GenerateRandomID() string {
	return uuid.New().String()
}
