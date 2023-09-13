package handler

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"

	"proto/hash_service"
)

func TestHashStrings(t *testing.T) {
	// Создаем тестовый сервер
	server := &Server{Logger: logrus.StandardLogger()}

	// Создаем тестовый клиент и передаем ему сервер
	client := NewTestClient(server)

	in := [][]byte{[]byte("hello"), []byte("world")}
	var expectedOut []string
	for _, str := range in {
		sha := sha3.New256()
		_, _ = sha.Write(str)
		res := hex.EncodeToString(sha.Sum(nil))
		expectedOut = append(expectedOut, res)
	}

	// Создаем тестовый запрос
	request := &hash_service.HashStringsRequest{
		Str: in,
	}

	// Отправляем запрос на сервер
	response, err := client.HashStrings(context.Background(), request)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Проверяем, что ответ содержит ожидаемые хэши
	assert.Equal(t, expectedOut, response.Hashed)
}

// TestClient реализует интерфейс HashServiceClient для тестирования
type TestClient struct {
	server *Server
}

// NewTestClient создает новый тестовый клиент
func NewTestClient(server *Server) hash_service.HashServiceClient {
	return &TestClient{server}
}

// HashStrings отправляет запрос на сервер
func (c *TestClient) HashStrings(ctx context.Context, in *hash_service.HashStringsRequest, opts ...grpc.CallOption) (*hash_service.HashStringsResponse, error) {
	return c.server.HashStrings(ctx, in)
}
