package config

type Config struct {
	Backends     []Backend    `yaml:"backends"`
	LoadBalancer LoadBalancer `yaml:"load-balancer"`
	Health       Health       `yaml:"healthcheck"`
}

type Backend struct {
	URL string `yaml:"url"`
}

type LoadBalancer struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Algorithm string `yaml:"algorithm"`
}

type Health struct {
	Path string `yaml:"path"`
}
