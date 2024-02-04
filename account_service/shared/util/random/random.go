package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet string = "abcdefghijklmnopqrstuvwxyz"
	integer  string = "1234567890"
)

var TransactionType = [2]string{"Buy", "Sell"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomFloat(max float64) float64 {
	return rand.Float64() * max
}

func RandomTransactionType() string {
	return TransactionType[rand.Intn(len(TransactionType))]
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomEmailString(n int) string {
	host := RandomString(n)
	return host + "@gmail.com"
}

func RandomStrInt(n int) string {
	return fmt.Sprintf("%d", rand.Intn(n))
}

func RandomFloatString(n int, point int) string {
	var sb strings.Builder
	k := len(integer)
	for i := 0; i < n; i++ {
		c := integer[rand.Intn(k)]
		if i == 0 && c == '0' {
			i--
			continue
		}
		sb.WriteByte(c)
	}
	segma := byte('.')
	sb.WriteByte(segma)
	for i := 0; i < point; i++ {
		c := integer[rand.Intn(k)]
		if i == 0 && c == '0' {
			i--
			continue
		}
		sb.WriteByte(c)
	}
	return sb.String()
}
