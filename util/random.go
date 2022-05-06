package util

import (
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxy"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// RandomInt generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomString generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomBalance generates a random balance
func RandomBalance() string {
	d, _ := decimal.NewFromString(strconv.Itoa(RandomInt(100, 200)))
	return d.StringFixed(4)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
