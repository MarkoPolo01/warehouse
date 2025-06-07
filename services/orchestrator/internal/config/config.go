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
				URL:  "http://localhost:8082/api/v1/abc-placement",
			},
			"fixed": {
				Name: "Fixed Placement",
				URL:  "http://localhost:8080/api/v1/fixed-placement",
			},
			"free": {
				Name: "Free Placement",
				URL:  "http://localhost:8081/process-placement",
			},
			"genetic": {
				Name: "Genetic Placement",
				URL:  "http://localhost:8085/api/v1/genetic-placement",
			},
			"greedy": {
				Name: "Greedy Placement",
				URL:  "http://localhost:8084/api/v1/greedy-placement",
			},
			"xyz": {
				Name: "XYZ Placement",
				URL:  "http://localhost:8083/api/v1/xyz-placement",
			},
		},
	}
} 