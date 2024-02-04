package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/api/gapi/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, subject *token.TokenSubject, audience string, issuer string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(
		subject,
		audience,
		issuer,
		duration,
	)
	require.NoError(t, err)
	brarerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, accessToken)
	md := metadata.MD{
		authorizationHeaderKey: []string{
			brarerToken,
		},
	}
	return metadata.NewIncomingContext(context.Background(), md)
}
