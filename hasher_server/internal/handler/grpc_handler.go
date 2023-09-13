package handler

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"hasher_server/pkg/hasher"
	"proto/hash_service"
)

type Server struct {
	Strings hash_service.HashStringsRequest
	hash_service.UnimplementedHashServiceServer
	Logger *logrus.Logger
}

func (s *Server) HashStrings(ctx context.Context, in *hash_service.HashStringsRequest) (*hash_service.HashStringsResponse, error) {
	s.Logger.Info("HashStrings called")
	strinsBytes := in.GetStr()
	var wg sync.WaitGroup
	hashes := make([]string, len(strinsBytes))

	for i, str := range strinsBytes {
		wg.Add(1)
		i := i
		go func(str string) {
			defer wg.Done()
			hash, err := hasher.GetHashedStr([]byte(str), s.Logger)
			if err != nil {
				s.Logger.WithFields(logrus.Fields{"StackTrace": errors.WithStack(err)}).Errorf("HashStrings error: %s", err)
				return
			}
			hashes[i] = *hash
		}(string(str))
	}
	wg.Wait()
	return &hash_service.HashStringsResponse{Hashed: hashes}, nil
}
