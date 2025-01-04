// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	dbstorage "github.com/fasdalf/train-go-gophkeeper/internal/server/db/storage"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
)

var typeCheckSecretDBRepository repository.SecretRepository = &SecretDBRepository{}

type SecretDBRepository struct {
	dp *dbstorage.DBProxy
	// TODO: ##@@ implement methods
}

// NewSecretDBRepository init
func NewSecretDBRepository(dp *dbstorage.DBProxy) *SecretDBRepository {
	return &SecretDBRepository{
		dp: dp,
	}
}

func (r *SecretDBRepository) FindById(Id uint64) (*entity.Secret, error) {
	return nil, nil
}
func (r *SecretDBRepository) FindAfter(UserId uint64, UpdatedAt uint64) ([]*entity.Secret, error) {

	// TODO: ##@@ #cleanup
	// ##@@ for multiple results use internal/common/metricstorage/dbStorage.go:160
	//for rows.Next() {
	//	var k string
	//	if err = rows.Scan(&k); err != nil {
	//		return nil, err
	//	}
	//	keys = append(keys, k)
	//}
	//
	//if err = rows.Err(); err != nil {
	//	return nil, err
	//}

	return nil, nil
}
func (r *SecretDBRepository) Create(UserId uint64, Data []byte) (*entity.Secret, error) {
	return nil, nil
}
func (r *SecretDBRepository) Update(Id uint64, Data []byte, OldUpdatedAt uint64) (*entity.Secret, error) {
	return nil, nil
}
func (r *SecretDBRepository) Delete(Id uint64, OldUpdatedAt uint64) (*entity.Secret, error) {
	return nil, nil
}
