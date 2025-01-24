package config

import "time"

type Config struct {
	Server *Server `yaml:"server"`
	DB     *DB     `yaml:"db"`
}

type Server struct {
	Port              int           `yaml:"port"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
}

type DB struct {
	ConnectionString string `yaml:"connectionString"`
}
