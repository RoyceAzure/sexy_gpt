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
type TokenSubject struct {
	ID     uuid.UUID `json:"id"`
	UPN    string    `json:"upn"`
	UserId uuid.UUID `json:"user_id"`
	RoleId uuid.UUID `json:"role_id"`
}

type TokenPayload struct {
	Audience  string        `json:"audience"`
	Issuer    string        `json:"issuer"`
	IssuedAt  time.Time     `json:"issuer_at"`
	ExpiredAt time.Time     `json:"expired_at"` // 短期有效
	Subject   *TokenSubject `json:"subject"`
}

func NewTokenSubject(upn string, userId uuid.UUID, roleId uuid.UUID) *TokenSubject {
	id := uuid.New()
	tokenSubject := &TokenSubject{
		ID:     id,
		UPN:    upn,
		UserId: userId,
		RoleId: roleId,
	}
	return tokenSubject
}

func NewTokenPayload(tokenSubject *TokenSubject, audience string, issuer string, duration time.Duration) *TokenPayload {
	now := time.Now().UTC()
	payload := &TokenPayload{
		Audience:  audience,
		Issuer:    issuer,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
		Subject:   tokenSubject,
	}
	return payload
}

// todo  驗證issueer, audience
func (payload *TokenPayload) Valid() error {
	if time.Now().UTC().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
