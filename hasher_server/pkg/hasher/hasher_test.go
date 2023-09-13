package hasher

import (
	"encoding/hex"
	"testing"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

func TestGetHashedStr(t *testing.T) {
	// Сценарий 1: Пустая строка
	emptyData := ""
	sha := sha3.New256()
	_, _ = sha.Write([]byte(emptyData))
	expectedEmptyHash := hex.EncodeToString(sha.Sum(nil))
	hashEmpty, err := GetHashedStr([]byte{}, logrus.StandardLogger())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if *hashEmpty != expectedEmptyHash {
		t.Fatalf("Expected hash %s for empty string, got %s", expectedEmptyHash, *hashEmpty)
	}

	// Сценарий 2: Строка с пробелами
	stringWithSpaces := "   hello   "
	sha = sha3.New256()
	_, _ = sha.Write([]byte(stringWithSpaces))
	expectedTrimmedHash := hex.EncodeToString(sha.Sum(nil))
	hashSpaces, err := GetHashedStr([]byte(stringWithSpaces), logrus.StandardLogger())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if *hashSpaces != expectedTrimmedHash {
		t.Fatalf("Expected hash %s for trimmed string, got %s", expectedTrimmedHash, *hashSpaces)
	}

	// Добавьте здесь еще сценарии по своему выбору
}
