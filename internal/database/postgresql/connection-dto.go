package postgresql

type ConnectionDTO struct {
	HostName string
	User     string
	Password string
	Port     string
	DBName   string
}
