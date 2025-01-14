package dbrepository

import (
	"context"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"testing"
)

func TestSecretMockRepository(t *testing.T) {
	se := &entity.Secret{
		ID:        10,
		UserId:    1,
		UpdatedAt: 500,
	}
	c := context.Background()
	r := SecretMockRepository{
		Secret:  se,
		Secrets: []*entity.Secret{se},
		Error:   nil,
	}
	r.FindById(c, 10)
	r.FindAfter(c, 1, 500)
	r.Create(c, 10, nil)
	r.Update(c, 10, nil, 500)
	r.Delete(c, 10, 500)
}
