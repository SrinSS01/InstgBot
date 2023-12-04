package info

import "time"

type Info struct {
	Session  string
	Username string
	Name     string
	Bio      string
	URL      string
	Email    string
	Phone    string
	Timer    *time.Timer
	Question int
	Posts    []string
	Pfp      string
}
