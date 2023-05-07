// Package user - this is a domain part of user that contains user' model and service functionality.
package user

type User struct {
	User       string  `json:"login"`
	Password   string  `json:"password"`
	Accrual    float64 `json:"current,omitempty"`
	Withdrawal float64 `json:"withdrawn,omitempty"`
}

func (u User) Equals(u2 *User) (ok bool) {
	if u2 == nil {
		return false
	}
	if u.User == u2.User && u.Password == u2.Password {
		return true
	}
	return false
}
