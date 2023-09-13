package hasher

import (
	"encoding/hex"

	er "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

func GetHashedStr(str []byte, logger *logrus.Logger) (*string, error) {

	sha := sha3.New256()
	_, err := sha.Write(str)
	if err != nil {
		logger.WithFields(logrus.Fields{"StackTrace": er.WithStack(err)}).Errorf("GetHashedStr error: %s", err)
		return nil, err
	}
	res := hex.EncodeToString(sha.Sum(nil))
	return &res, nil
}
