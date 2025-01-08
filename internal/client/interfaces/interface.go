// Package interfaces - break circle dependencies.
// TODO: ##@@ extract to using packages
package interfaces

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
)

type RemoteService interface {
	SignIn(login, pass string) (string, error)
	SignUp(login, pass string) (string, error)
	SetToken(token string)
	ListSecrets(fromTs int64) (*pb.ListSecretsResponse, error)
	CreateSecret(data *[]byte) (id, ts int64, err error)
	UpdateSecret(id, oldTs int64, data *[]byte) (ts int64, err error)
	DeleteSecret(id, oldTs int64) (err error)
}

type Syncroniser interface {
	Sync() error
	Add(secret entity.Secret) (serverID int64, serverUpdatedAt int64, err error)
	Update(secret entity.Secret) (serverID int64, serverUpdatedAt int64, err error)
	Delete(localId int) error
}

type LocalRepository interface {
	ListSecrets() []string
	ClearSecrets()
	AddSecret(secret entity.Secret) (int, error)
	GetSecret(id int) (entity.Secret, error)
	GetSecretIdByServerId(id int64) (int, error)
	SetSecret(id int, secret entity.Secret) (entity.Secret, error)
	DeleteSecret(id int) error
}

type ButtonHandler interface {
	Handle(m tea.Model) tea.Model
}
