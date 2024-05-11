package config

type PoolConfig struct {
	MaxIdleConns int `json:"max_idle_conns" envconfig:"DB_POOL_MAXIDLECONNS"`
	MaxOpenConns int `json:"max_open_conns" envconfig:"DB_POOL_MAXOPENCONNS"`
	IdleTimeout  int `json:"idle_timeout" envconfig:"IDLE_TIMEOUT"`
}

type DBConfig struct {
	Host     string `json:"host" envconfig:"DB_HOST"`
	Port     string `json:"port" envconfig:"DB_PORT"`
	Database string `json:"database" envconfig:"DB_DATABASE"`
	Username string `json:"username" envconfig:"DB_USERNAME"`
	Password string `json:"password" envconfig:"DB_PASSWORD"`
	Pool     PoolConfig
}