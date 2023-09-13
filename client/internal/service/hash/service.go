package hash_service

import (
	"context"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/pkg/errors"

	"client/internal/handler/operations"
	"client/internal/repository/hash_repository"
	"client/models"
	"hasher_server/pkg/pb/hash_service"
)

type HashWorker interface {
	HashStrings(ctx context.Context, in *hash_service.HashStringsRequest, opts ...grpc.CallOption) (*hash_service.HashStringsResponse, error)
}

type HashService struct {
	repository *hash_repository.HashRepo
	hw         HashWorker
	logger     *logrus.Logger
}

func New(repository *hash_repository.HashRepo, hw HashWorker, l *logrus.Logger) *HashService {
	return &HashService{repository: repository, hw: hw, logger: l}
}

func (s *HashService) GetHashedString(p *operations.GetCheckParams) middleware.Responder {
	s.logger.WithFields(logrus.Fields{"requestID": p.HTTPRequest.Context().Value("X-Request-Id")}).Info("GetHashedString called")
	var res []*models.Hash
	for _, v := range p.Ids {
		id, err := strconv.Atoi(v)
		if err != nil {
			s.logger.WithFields(logrus.Fields{"StackTrace": errors.WithStack(err), "requestID": p.HTTPRequest.Context().Value("X-Request-Id")}).Errorf("invalid id: %s, error: %s", v, err)
			return operations.NewPostSendBadRequest()
		}
		hashFromDB, err := s.repository.GetHashedString(id)
		if err != nil {
			s.logger.WithFields(logrus.Fields{"StackTrace": errors.WithStack(err), "requestID": p.HTTPRequest.Context().Value("X-Request-Id")}).Errorf("GetHashedString error: %s", errors.WithStack(err))
		}
		if hashFromDB != nil {
			hash := &models.Hash{
				Hash: &hashFromDB.Hash,
				ID:   &hashFromDB.ID,
			}
			res = append(res, hash)
		}
	}
	return operations.NewPostSendOK().WithPayload(res)
}

func (s *HashService) WriteHashedString(p *operations.PostSendParams) middleware.Responder {
	s.logger.WithFields(logrus.Fields{"requestID": p.HTTPRequest.Context().Value("X-Request-Id")}).Info("WriteHashedString called")
	if len(p.Params) == 0 {
		return operations.NewPostSendBadRequest()
	}
	var inputBytes [][]byte
	for _, str := range p.Params {
		inputBytes = append(inputBytes, []byte(str))
	}
	res, err := s.hw.HashStrings(context.Background(), &hash_service.HashStringsRequest{Str: inputBytes})
	if err != nil {
		s.logger.WithFields(logrus.Fields{"StackTrace": errors.WithStack(err), "requestID": p.HTTPRequest.Context().Value("X-Request-Id")}).Errorf("HashStrings error: %s", err)
		return operations.NewPostSendInternalServerError()
	}
	output := make(models.ArrayOfHash, len(res.Hashed))
	for i, r := range res.Hashed {
		hashed, err := s.repository.WriteHashedString(&models.Hash{Hash: &r})
		if err != nil {
			s.logger.WithFields(logrus.Fields{"StackTrace": errors.WithStack(err), "requestID": p.HTTPRequest.Context().Value("X-Request-Id")}).Errorf("WriteHashedString error: %s", err)
			return operations.NewPostSendInternalServerError()
		}
		output[i] = &models.Hash{Hash: &hashed.Hash, ID: &hashed.ID}
	}

	return operations.NewGetCheckOK().WithPayload(output)
}
