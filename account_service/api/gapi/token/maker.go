package token

import (
	"time"
)

type Maker interface {
	CreateToken(subject *TokenSubject, audience string, issuer string, duration time.Duration) (string, *TokenPayload, error)
	/*
		驗證簽章  過期 audi, iss
	*/
	VertifyToken(token string) (*TokenPayload, error)
}
