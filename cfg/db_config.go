package cfg

type DBConfig struct {
	User     string
	Password string
	Name     string
}

var DBCfg = DBConfig{
	User:     "root",
	Password: "", 
	Name:     "sapper",
}
