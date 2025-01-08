// Package localrepository - local repository
package localrepository

import (
	"errors"
	ifcs "github.com/fasdalf/train-go-gophkeeper/internal/client/interfaces"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
)

var typeCheckLocalRepository ifcs.LocalRepository = &LocalRepository{}
var ErrOutOfBounds = errors.New("out of bounds")

type localItem struct {
	name   string
	secret entity.Secret
}
type LocalRepository struct {
	items      []localItem
	byServerId map[int64]int
}

func NewLocalRepository() *LocalRepository {
	return &LocalRepository{
		items:      make([]localItem, 0),
		byServerId: make(map[int64]int),
	}
}

func (r *LocalRepository) ListSecrets() []string {
	result := make([]string, len(r.items))
	for i := range r.items {
		result[i] = r.items[i].name
	}
	return result
}

func (r *LocalRepository) ClearSecrets() {
	r.items = make([]localItem, 0)
	r.byServerId = make(map[int64]int)
}

func (r *LocalRepository) AddSecret(secret entity.Secret) (int, error) {
	i := len(r.items)
	data, _ := secret.GetData()
	li := localItem{
		name:   data.Name,
		secret: secret,
	}
	r.items = append(r.items, li)
	r.byServerId[secret.ServerID] = i

	return i, nil
}

func (r *LocalRepository) GetSecret(id int) (entity.Secret, error) {
	if id < 0 || id >= len(r.items) {
		return entity.Secret{}, ErrOutOfBounds
	}
	return r.items[id].secret, nil
}

func (r *LocalRepository) GetSecretIdByServerId(serverId int64) (int, error) {
	id, ok := r.byServerId[serverId]
	if !ok || id < 0 || id >= len(r.items) {
		return 0, ErrOutOfBounds
	}
	return id, nil
}

func (r *LocalRepository) SetSecret(id int, secret entity.Secret) (entity.Secret, error) {
	if id < 0 || id >= len(r.items) {
		return entity.Secret{}, ErrOutOfBounds
	}
	data, _ := secret.GetData()
	li := localItem{
		name:   data.Name,
		secret: secret,
	}
	r.items[id] = li
	return secret, nil
}

func (r *LocalRepository) DeleteSecret(id int) error {
	if id < 0 || id >= len(r.items) {
		return ErrOutOfBounds
	}
	li := r.items[id]
	delete(r.byServerId, li.secret.ServerID)

	r.items = append(r.items[:id], r.items[id+1:]...)
	for i := id; i < len(r.items); i++ {
		li = r.items[i]
		r.byServerId[li.secret.ServerID] = i
	}

	return nil
}
