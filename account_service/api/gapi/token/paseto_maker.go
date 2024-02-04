package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symmerickey []byte
}

func NewPasetoMaker(symmerickey string) (Maker, error) {
	if len(symmerickey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid ket size : must be exactly %d charcters", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{paseto.NewV2(), []byte(symmerickey)}, nil
}

// err: encode err
func (maker *PasetoMaker) CreateToken(subject *TokenSubject, audience string, issuer string, duration time.Duration) (string, *TokenPayload, error) {
	tokenPayload := NewTokenPayload(subject, audience, issuer, duration)
	token, err := maker.paseto.Encrypt(maker.symmerickey, tokenPayload, nil)
	return token, tokenPayload, err
}

func (maker *PasetoMaker) VertifyToken(token string) (*TokenPayload, error) {
	payload := &TokenPayload{}
	err := maker.paseto.Decrypt(token, maker.symmerickey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
