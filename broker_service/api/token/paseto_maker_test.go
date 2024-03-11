package token

import (
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/broker_service/shared/util/random"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker("12345678123456781234567812345678")
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := random.RandomString(10)
	userID := random.RandomInt(100, 1000)
	duration := time.Minute

	issuedAt := time.Now().UTC()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, userID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VertifyToken("v2.local.RlMsab1TWTJp5Gr61g51MNnZI0SinBTgO2k8HascdiGhdjj2OdzDWfdjjnXr7F6Ez0eZO05A2GHKBiQN_J0lnDWVcw91U4vx2vXh2ZXYB0sWzwW3h1A76CaElONUpW3CAMrm6-LhOqdfG-m0SQ7B44RqCtlC_3TS3M9B8gX6hdIyLpdkMPVuRghAzUobEJhw4UTAO0bbn2SGIaSZzbnOwKbFJD27_zvzibfdDxY5675uHdRMEMBekB43jp2voG4L0tUuz6V4RRlryPTBbo19P4B7phZye1wUZzVi_N7LfR5cZSckiVBAnaiOYA3NXdpptMRmoKCEZYC95UU29tfew9k1joBHHk4V3Dk89T8hNIksJpDD1XlxCtehDMfg4mKcQi0AdbjnBvz9ItPBdOBclfeEBKrFQpq69q8RXooG7NQct3XgiYwi0HR5vROft5qaCMay_N6yX_Y.bnVsbA")
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.UPN)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := random.RandomString(10)
	userID := random.RandomInt(100, 1000)
	duration := time.Minute

	//createToken沒有禁止使用負數的duration
	token, payload, err := maker.CreateToken(username, userID, -duration)
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
	payload, err := NewPayload(random.RandomString(10), 1, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
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
