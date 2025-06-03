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
			"fixed": {
				Name: "Fixed Placement",
				URL:  "http://localhost:8080",
			},
			"abc": {
				Name: "ABC Placement",
				URL:  "http://localhost:8081",
			},
			"xyz": {
				Name: "XYZ Placement",
				URL:  "http://localhost:8082",
			},
			"dynamic": {
				Name: "Dynamic Placement",
				URL:  "http://localhost:8083",
			},
			"optimal": {
				Name: "Optimal Placement",
				URL:  "http://localhost:8084",
			},
		},
	}
} 