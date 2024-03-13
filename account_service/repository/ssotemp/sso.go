package sso

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/square/go-jose/v3/jwt"
	"golang.org/x/crypto/hkdf"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

type GoogleSSO struct {
	jwt.Claims
	Nmae        string `json:"name"`
	Email       string `json:"email"`
	IDToken     string `json:"id_token"`
	AccessToken string `json:"accessToken"`
}

// use hkdf and sha256 encrypted key, then decode encryptedToken
// validate idtoken in encryptedToken used idtoken.NewValidator
func VertifyGoogleSSOIDToken(ctx context.Context, encryptedToken, key string) (*GoogleSSO, error) {
	encryped_key, err := encryptionKey(key)
	if err != nil {
		return nil, err
	}

	googleSSO, err := decodeJWE(encryptedToken, encryped_key)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	validator, err := idtoken.NewValidator(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error create idtoken validator : %s, %w", err.Error(), gpt_error.ErrInternal)
	}

	_, err = validator.Validate(ctx, googleSSO.IDToken, "")
	if err != nil {
		return nil, fmt.Errorf(" %s, %w", err.Error(), gpt_error.ErrUnauthicated)
	}

	return googleSSO, nil
}

// use hkdf and sha256 to encryped secret
func encryptionKey(secret string) ([]byte, error) {
	salt := make([]byte, 0)
	hkdf := hkdf.New(sha256.New, []byte(secret), salt, []byte("NextAuth.js Generated Encryption Key"))
	key := make([]byte, 32)
	_, err := hkdf.Read(key)
	if err != nil {
		return salt, fmt.Errorf(" failed encrypt key : %s, %w", err.Error(), gpt_error.ErrInternal)
	}
	return key, nil
}

// decode token send from frontend, need secret
// use sha256 and hkdf to decode secret
// return GoogleSSO struct if successed
func decodeJWE(token string, secret []byte) (*GoogleSSO, error) {
	obj, err := jwt.ParseEncrypted(token)
	if err != nil {
		return nil, fmt.Errorf("decode JWE failed,%s,  %w", err.Error(), gpt_error.ErrInValidatePreConditionOp)
	}

	out := GoogleSSO{}
	if err := obj.Claims(secret, &out); err != nil {
		return nil, fmt.Errorf("failed to read google sso claims,%s,  %w", err.Error(), gpt_error.ErrInternal)
	}

	return &out, nil
}
