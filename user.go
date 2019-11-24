package gozk

type User struct {
	UserID     int
	UserSN     string
	Name       string
	Password   string
	CardNo     string
	AdminLevel int
	Enabled    bool
	GroupNo    int
	Timezones  []string
}
