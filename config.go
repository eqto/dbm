package dbm

type Config struct {
	DriverName string

	Hostname string
	Port     int
	Username string
	Password string
	Name     string

	DisableEncryption bool
}
