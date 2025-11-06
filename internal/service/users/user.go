package service

import (
	"productsParser/internal/domain"
	i "productsParser/internal/domain/interface"
	service "productsParser/internal/service/orders"
	"sync"
)

type UserService struct {
	users map[string]*service.Orders
	pInit func() i.Parser
	rw    *sync.RWMutex
}

func NewUserService(p func() i.Parser) *UserService {
	return &UserService{
		users: make(map[string]*service.Orders),
		pInit: p,
		rw:    &sync.RWMutex{},
	}
}

func (u *UserService) GetOrAddUserService(id string) (*service.Orders, bool) {
	user, ok := u.GetUserService(id)

	if !ok {
		u.rw.Lock()
		defer u.rw.Unlock()
		us := service.NewOrderService(domain.NewUser(id), u.pInit())
		user, ok := u.users[id]
		if !ok {
			u.users[id] = us
			return us, ok
		}
		return user, ok
	}

	return user, ok
}

func (u *UserService) GetUserService(id string) (*service.Orders, bool) {
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
