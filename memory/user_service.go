package memory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cloudfoundry-community/go-uaa"

	"codeberg.org/ess/fuaa/core"
)

type UserService struct {
	users  []uaa.User
	tokens *TokenService
	locker sync.Mutex
}

func NewUserService(tokens *TokenService) *UserService {
	service := &UserService{tokens: tokens}
	service.Reset()

	return service
}

func (service *UserService) ByUsername(username string) (uaa.User, error) {
	service.locker.Lock()
	defer service.locker.Unlock()

	return service.byUsername(username)
}

func (service *UserService) byUsername(username string) (uaa.User, error) {
	fmt.Println("got username:", username)
	for _, user := range service.users {
		fmt.Println("Comparing", user.Username, "to", username)
		if user.Username == username {
			return user, nil
		}
	}

	return uaa.User{}, errors.New("user not found")
}

func (service *UserService) Add(username string, password string) (uaa.User, error) {
	service.locker.Lock()
	defer service.locker.Unlock()

	return service.add(username, password)
}

func (service *UserService) add(username string, password string) (uaa.User, error) {
	if service.exists(username) {
		return uaa.User{}, errors.New("already exists")
	}

	userGUID, _ := core.GenGUID()

	user := uaa.User{
		ID:       userGUID,
		Username: username,
		Password: password,
		Origin:   "here",
	}

	_, err := service.tokens.Create(user)
	if err != nil {
		return uaa.User{}, errors.New("could not create token for user")
	}

	service.users = append(service.users, user)

	return user, nil
}

func (service *UserService) exists(username string) bool {
	for _, user := range service.users {
		if user.Username == username {
			return true
		}
	}

	return false
}

func (service *UserService) Reset() {
	service.locker.Lock()
	defer service.locker.Unlock()

	service.reset()
}

func (service *UserService) reset() {
	service.users = make([]uaa.User, 0)
	service.add("admin", "admin")
}
