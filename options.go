package dbm

import "time"

type Options func(*Connection)

//OptionMaxIdleTime default is 60 seconds
func OptionMaxIdleTime(duration time.Duration) Options {
	return func(c *Connection) {
		c.db.SetConnMaxIdleTime(duration)
	}
}

//OptionMaxLifetime default is 60 minutes
func OptionMaxLifetime(duration time.Duration) Options {
	return func(c *Connection) {
		c.db.SetConnMaxLifetime(duration)
	}
}

//OptionMaxIdle default is 2
func OptionMaxIdle(count int) Options {
	return func(c *Connection) {
		c.db.SetMaxIdleConns(count)
	}
}

//OptionMaxOpen default is 50
func OptionMaxOpen(count int) Options {
	return func(c *Connection) {
		c.db.SetMaxOpenConns(count)
	}
}

func OptionDisableEncryption() Options {
	return func(c *Connection) {
		c.cfg.DisableEncryption = true
	}
}
