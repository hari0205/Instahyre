package ini

import (
	"os"

	model "example.com/Instahyre/teleapi/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	err_env := godotenv.Load()
	if err_env != nil {
		panic("Error loading environment variables")
	}
	var err error
	dsn := os.Getenv("DB_URL") //"host=tiny.db.elephantsql.com user=mpgqlmdw password=qsq_2KLu5wlBRcPtr2karRdIpiiwTLpD dbname=mpgqlmdw port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	DB.AutoMigrate(&model.UserData{})
}
