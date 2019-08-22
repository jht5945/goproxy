package config

import (
	//"encoding/json"
)

type HostMapping struct {
	Host string
	Port string
	Target string
}

type Config struct {
	ListenPorts []string
	HostMappings []HostMapping
}

/*
func LoadConfig() Config {
	json.Unmarshal()
}
*/
