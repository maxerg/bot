package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
	"time"
)

type InitData struct {
	QueryID  string
	UserJSON string
	AuthDate int64
	Hash     string

	// extra fields may exist
	Raw map[string]string
}

func ValidateInitData(initData, botToken string, maxAge time.Duration) (InitData, bool) {
	initData = strings.TrimSpace(initData)
	botToken = strings.TrimSpace(botToken)
	if initData == "" || botToken == "" {
		return InitData{}, false
	}

	v, err := url.ParseQuery(initData)
	if err != nil {
		return InitData{}, false
	}

	hash := v.Get("hash")
	if hash == "" {
		return InitData{}, false
	}

	// Build data-check-string: key=<value> lines sorted by key
	keys := make([]string, 0, len(v))
	for k := range v {
		// IMPORTANT: exclude both hash and signature (signature is for 3rd-party validation)
		if k == "hash" || k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	raw := make(map[string]string, len(keys))
	for _, k := range keys {
		val := v.Get(k) // decoded value (as in Telegram docs examples)
		raw[k] = val
		lines = append(lines, k+"="+val)
	}
	dataCheckString := strings.Join(lines, "\n")

	// secret_key = HMAC_SHA256(bot_token, "WebAppData")  (key="WebAppData", msg=bot_token)
	m1 := hmac.New(sha256.New, []byte("WebAppData"))
	m1.Write([]byte(botToken))
	secretKey := m1.Sum(nil)

	// computed_hash = hex(HMAC_SHA256(data_check_string, secret_key))
	m2 := hmac.New(sha256.New, secretKey)
	m2.Write([]byte(dataCheckString))
	computed := hex.EncodeToString(m2.Sum(nil))

	if !hmac.Equal([]byte(computed), []byte(hash)) {
		return InitData{}, false
	}

	// Optional anti-replay by auth_date
	authDateStr := v.Get("auth_date")
	authDate, _ := parseInt64(authDateStr)
	if authDate == 0 {
		return InitData{}, false
	}
	if maxAge > 0 {
		now := time.Now().Unix()
		if now-authDate > int64(maxAge.Seconds()) {
			return InitData{}, false
		}
	}

	out := InitData{
		QueryID:  v.Get("query_id"),
		UserJSON: v.Get("user"),
		AuthDate: authDate,
		Hash:     hash,
		Raw:      raw,
	}
	return out, true
}

func parseInt64(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	// fast path: decimal only
	var n int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, nil
		}
		n = n*10 + int64(c-'0')
	}
	return n, nil
}
