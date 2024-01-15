package db

import (
	"crypto/tls"
	"fmt"
	"time"
)

type Config struct {
	Host              string        `json:"host" yaml:"host"`
	Port              uint16        `json:"port" yaml:"port"`
	Database          string        `json:"database" yaml:"database"`
	User              string        `json:"user" yaml:"user"`
	Password          string        `json:"password" yaml:"password"`
	TLSConfig         *tls.Config   `json:"tls_config" yaml:"tls_config"` // nil disables TLS
	ConnectTimeout    time.Duration `json:"connect_timeout" yaml:"connect_timeout"`
	MaxConn           int64         `json:"max_conn" yaml:"max_conn"`
	MinConn           int64         `json:"min_conn" yaml:"min_conn"`
	HealCheckDuration time.Duration `json:"heal_check_duration" yaml:"heal_check_duration"`
}

func (db *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
	)
}
