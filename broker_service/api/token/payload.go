package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// 這個Payload也等同於Claim   這個套件的Valid完全由自己掌控??  不對  只有claim需要自己驗證  其餘簽名應由套件處理
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UPN       string    `json:"upn"`
	UserId    int64     `json:"userid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `josn:"expired_at"`
}

func NewPayload(upn string, userID int64, duration time.Duration) (*Payload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        uuid,
		UPN:       upn,
		UserId:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// 需要實現jwt Claim的Valid街口  反正就是你的claim資料要自己寫驗證
func (payload *Payload) Valid() error {
	if time.Now().UTC().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
