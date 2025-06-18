package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ascii = `Testing load balancer concept.
`
)

type Config struct {
	SSLCertificateKey string   `yaml:"ssl_certificate_key"`
	ServersPath       []string `yaml:"servers_path"`
	Port              int      `yaml:"lb_port"`
	HealthCheck       bool     `yaml:"tcp_health_check"`
	HealthCheckPath   string   `yaml:"health_check_path"`
	MaxAllowed        uint     `yaml:"max_allowed"`
}

func ReadConfig(fileName string) (*Config, error) {
	in, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Print() {
	fmt.Printf("%s\nPort: %d\nHealth Check: %v\nLocation:\n",
		ascii, c.Port, c.HealthCheck)
	// for _, l := range c.Location {
	// 	fmt.Printf("\tRoute: %s\n\tProxy Pass: %s\n\tMode: %s\n\n",
	// 		l.Pattern, l.ProxyPass, l.BalanceMode)
	// }
}
