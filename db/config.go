package db

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

type PostgresConnectionConfig struct {
	Host   string
	Port   int
	User   string
	Passwd string
	DBName string
}

func (c *PostgresConnectionConfig) DSN() string {
	userEsc := url.QueryEscape(c.User)
	passEsc := url.QueryEscape(c.Passwd)
	hostPort := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", userEsc, passEsc, hostPort, c.DBName)

}
