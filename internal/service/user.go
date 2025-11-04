package service

import (
	"productsParser/internal/domain"
	"sync"
)

type UserService struct {
	users map[string]*OrdersService
	pInit func() Parser
	rw    *sync.RWMutex
}

func NewUserService(p func() Parser) *UserService {
	return &UserService{
		users: make(map[string]*OrdersService),
		pInit: p,
		rw:    &sync.RWMutex{},
	}
}

func (u *UserService) GetOrAddUserService(id string) (*OrdersService, bool) {
	user, ok := u.GetUserService(id)

	if !ok {
		u.rw.Lock()
		defer u.rw.Unlock()
		us := NewOrderService(domain.NewUser(id), u.pInit())
		user, ok := u.users[id]
		if !ok {
			u.users[id] = us
			return us, ok
		}
		return user, ok
	}

	return user, ok
}

func (u *UserService) GetUserService(id string) (*OrdersService, bool) {
	u.rw.RLock()
	defer u.rw.RUnlock()
	user, ok := u.users[id]
	return user, ok
}

func (u *UserService) DelUser(id string) {
	u.rw.Lock()
	defer u.rw.Unlock()
	delete(u.users, id)
}
