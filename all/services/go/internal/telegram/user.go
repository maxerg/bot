package telegram

import "encoding/json"

type WebAppUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func ParseUserJSON(s string) (WebAppUser, bool) {
	if s == "" {
		return WebAppUser{}, false
	}
	var u WebAppUser
	if err := json.Unmarshal([]byte(s), &u); err != nil {
		return WebAppUser{}, false
	}
	if u.ID == 0 {
		return WebAppUser{}, false
	}
	return u, true
}
