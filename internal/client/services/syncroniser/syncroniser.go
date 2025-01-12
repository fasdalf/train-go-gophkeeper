package syncroniser

import (
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
)

var typeCheckSyncroniser SyncroniserInterface = &Syncroniser{}

type SyncroniserInterface interface {
	Sync() error
	Add(secret entity.Secret) (serverID int64, serverUpdatedAt int64, err error)
	Update(secret entity.Secret) (serverID int64, serverUpdatedAt int64, err error)
	Delete(localId int) error
}

type LocalRepository interface {
	AddSecret(secret entity.Secret) (int, error)
	GetSecret(id int) (entity.Secret, error)
	GetSecretIdByServerId(id int64) (int, error)
	SetSecret(id int, secret entity.Secret) (entity.Secret, error)
	DeleteSecret(id int) error
}

type RemoteService interface {
	ListSecrets(fromTs int64) (*pb.ListSecretsResponse, error)
	CreateSecret(data *[]byte) (id, ts int64, err error)
	UpdateSecret(id, oldTs int64, data *[]byte) (ts int64, err error)
	DeleteSecret(id, oldTs int64) (err error)
}

type Syncroniser struct {
	Local      LocalRepository
	Remote     RemoteService
	SyncedUpTo int64
}

func (s *Syncroniser) Sync() error {
	remoteSecrets, err := s.Remote.ListSecrets(s.SyncedUpTo)
	if err != nil {
		return err
	}
	for _, rs := range remoteSecrets.GetSecrets() {
		lid, err := s.Local.GetSecretIdByServerId(rs.Id)
		if err != nil {
			// it's new, add
			ls := entity.Secret{
				ServerID:        rs.Id,
				ServerUpdatedAt: rs.Ts,
				Data:            &rs.Data,
			}
			_, _ = s.Local.AddSecret(ls)
			continue
		}

		ls, _ := s.Local.GetSecret(lid)

		if rs.IsDeleted {
			_ = s.Local.DeleteSecret(lid)
			continue
		}

		_, _ = s.Local.SetSecret(lid, ls)
	}

	s.SyncedUpTo = remoteSecrets.SyncTs
	return nil
}

func (s *Syncroniser) Add(secret entity.Secret) (serverID int64, serverUpdatedAt int64, err error) {
	serverID, serverUpdatedAt, err = s.Remote.CreateSecret(secret.Data)
	if err != nil {
		return
	}

	secret.ServerID = serverID
	secret.ServerUpdatedAt = serverUpdatedAt
	_, _ = s.Local.AddSecret(secret)
	return
}

func (s *Syncroniser) Update(secret entity.Secret) (serverID int64, serverUpdatedAt int64, err error) {
	serverID = secret.ServerID
	serverUpdatedAt, err = s.Remote.UpdateSecret(secret.ServerID, secret.ServerUpdatedAt, secret.Data)
	if err != nil {
		return
	}

	lid, _ := s.Local.GetSecretIdByServerId(serverID)
	_, _ = s.Local.SetSecret(lid, secret)
	return
}

func (s *Syncroniser) Delete(localId int) error {
	ls, err := s.Local.GetSecret(localId)
	if err == nil {
		err = s.Remote.DeleteSecret(ls.ServerID, ls.ServerUpdatedAt)
	}
	if err == nil {
		_ = s.Local.DeleteSecret(localId)
	}
	return err
}
