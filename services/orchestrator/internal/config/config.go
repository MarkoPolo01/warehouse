package config

// ServiceConfig содержит конфигурацию микросервиса
type ServiceConfig struct {
	Name string
	URL  string
}

// Config содержит конфигурацию всех микросервисов
type Config struct {
	Services map[string]ServiceConfig
}

// NewConfig создает новую конфигурацию
func NewConfig() *Config {
	return &Config{
		Services: map[string]ServiceConfig{
			"abc": {
				Name: "ABC Placement",
				URL:  "http://localhost:8082",
			},
			"fixed": {
				Name: "Fixed Placement",
				URL:  "http://localhost:8080",
			},
			"free": {
				Name: "Free Placement",
				URL:  "http://localhost:8081",
			},
			"genetic": {
				Name: "Genetic Placement",
				URL:  "http://localhost:8085",
			},
			"greedy": {
				Name: "Greedy Placement",
				URL:  "http://localhost:8084",
			},
			"xyz": {
				Name: "XYZ Placement",
				URL:  "http://localhost:8083",
			},
		},
	}
} 