package memory

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"sync"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/cloudfoundry-community/go-uaa"
)

type TokenService struct {
	keypair *rsa.PrivateKey
	tokens  map[string]string
	locker  sync.Mutex
}

func NewTokenService() *TokenService {
	service := &TokenService{}
	service.Reset()

	return service
}

func (service *TokenService) Create(user uaa.User) (string, error) {
	service.locker.Lock()
	defer service.locker.Unlock()

	return service.create(user)
}

func (service *TokenService) create(user uaa.User) (string, error) {
	jwt := jws.NewJWT(
		jws.Claims{
			"flibberty": "gibbets",
			"user_name": user.Username,
			"user_id":   user.ID,
			"origin":    "here",
		},
		crypto.SigningMethodRS256,
	)

	token, err := jwt.Serialize(service.keypair)
	if err != nil {
		return "", err
	}

	output := string(token)

	service.tokens[user.Username] = output

	return output, nil
}

func (service *TokenService) ByUser(user uaa.User) (string, error) {
	service.locker.Lock()
	defer service.locker.Unlock()

	return service.byUser(user)
}

func (service *TokenService) byUser(user uaa.User) (string, error) {
	token, ok := service.tokens[user.Username]
	if !ok {
		return "", errors.New("token not found")
	}

	return token, nil
}

func (service *TokenService) Exists(token string) bool {
	service.locker.Lock()
	defer service.locker.Unlock()

	return service.exists(token)
}

func (service *TokenService) exists(token string) bool {
	for _, candidate := range service.tokens {
		if candidate == token {
			return true
		}
	}

	return false
}

func (service *TokenService) Reset() {
	service.locker.Lock()
	defer service.locker.Unlock()

	service.reset()
}

func (service *TokenService) reset() {
	kp, _ := rsa.GenerateKey(rand.Reader, 2048)

	service.keypair = kp
	service.tokens = make(map[string]string)
}
