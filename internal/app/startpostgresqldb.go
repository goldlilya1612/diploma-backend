package app

import "github.com/goldlilya1612/diploma-backend/internal/database"

func StartPostgreSQL(config *database.Config) error {

	psql := database.NewPostgreSQL(config)
	err := psql.Connect()
	if err != nil {
		return err
	}

	err = psql.InitDB()
	if err != nil {
		return err
	}

	return nil
}
