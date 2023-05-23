package config

type Application struct {
	Address string
	Port    int
}

type Database struct {
	Name string
}

type GRPC struct {
	Address string
	Port    int
}

type Queue struct {
	Address  string
	Port     int
	Username string
	Password string
}
