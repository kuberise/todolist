package postgres

type Config struct {
	Host     string `cfg:"host"`
	Port     int    `cfg:"port"`
	User     string `cfg:"user"`
	Password string `cfg:"password"`
	DBName   string `cfg:"db_name"`
	SSLMode  string `cfg:"ssl_mode"`
}
