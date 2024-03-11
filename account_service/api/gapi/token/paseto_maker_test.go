package token

import (
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/random"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	// maker, err := NewPasetoMaker(random.RandomString(32))
	maker, err := NewPasetoMaker("12345678123456781234567812345678")
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	duration := time.Minute * 10

	issuedAt := time.Now().UTC()
	expiredAt := issuedAt.Add(duration)
	id := uuid.New()
	audi := "audi"
	isur := "isue"
	email := "roycewnag@gmail.com"
	subject := TokenSubject{
		ID:     id,
		UPN:    email,
		UserId: id,
		RoleId: id,
	}

	token, payload, err := maker.CreateToken(&subject, audi, isur, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VertifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, audi, payload.Audience)
	require.Equal(t, isur, payload.Issuer)
	require.Equal(t, email, payload.Subject.UPN)
	require.Equal(t, id, payload.Subject.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	duration := time.Minute

	audi := "audi"
	isur := "isue"

	token, payload, err := maker.CreateToken(nil, audi, isur, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VertifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoTokenAlgNone(t *testing.T) {
	//自己產生jwt token
	//payload
	payload := NewTokenPayload(nil, random.RandomString(10), random.RandomString(10), time.Minute)

	//選擇加密演算法製作claim
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	//指定key作加密
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err = maker.VertifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
