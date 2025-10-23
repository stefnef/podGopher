package env

import "os"
import "github.com/joho/godotenv"

type Name string

func (key Name) GetValue() string {
	return os.Getenv(string(key))
}

func Load(filename string) error {
	return godotenv.Load(filename)
}

const (
	DBName       Name = "DBName"
	DBUser       Name = "DBUser"
	DBPassword   Name = "DBPassword"
	DBHost       Name = "DBHost"
	DBPort       Name = "DBPort"
	MigrationDir Name = "MigrationDir"
)
