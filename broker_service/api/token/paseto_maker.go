package token

import (
	"fmt"
	"sync"
	"time"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/gpt_error"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symmerickey []byte
}

var (
	instance Maker
	once     sync.Once
)

func GetSingleTonPasetoTokenMaker() (Maker, error) {
	config, err := config.LoadConfig(".")
	if err != nil {
		return nil, fmt.Errorf("get singleTon token maker ger err : %s, %w", err.Error(), gpt_error.ErrInternal)
	}
	if instance == nil {
		once.Do(func() {
			instance, err = NewPasetoMaker(config.TokenSymmetricKey)
		})
	}
	return instance, nil
}

func NewPasetoMaker(symmerickey string) (Maker, error) {
	if len(symmerickey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid ket size : must be exactly %d charcters", chacha20poly1305.KeySize)
	}
	return &PasetoMaker{paseto.NewV2(), []byte(symmerickey)}, nil
}

func (maker *PasetoMaker) CreateToken(upn string, userID int64, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(upn, userID, duration)
	if err != nil {
		return "", payload, err
	}
	//相比jwt之所以只有這行是因為你不需要決定加密演算法
	//固定使用chacha演算法
	token, err := maker.paseto.Encrypt(maker.symmerickey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VertifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	//之所以這麼簡單  是因為套件只會回傳ErrInvalidToken相關範圍得錯誤
	//你自己的paload valid要自己呼叫
	//也因為如此  不像jwt驗證是通通包再一起，你必須拆解jwt回傳的錯誤訊息
	//目前看到的InvalidToken錯誤包括key  資料不匹配
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
