package mysql

import (
	"backend/internal/db"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestName(t *testing.T) {

	container := Container{}
	config := Configurations{
		Authentication: &db.AuthenticationConfigurations{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "password",
		},
		DatabaseName: "mydatabase",
	}

	err := container.Connect(config)

	_ = "SELECT uuid,username, hashed_pass, salt FROM UserTable WHERE username = ?"

	repo, err := CreateRepository(container, "CREATE TABLE IF NOT EXISTS UserTable(uuid VARCHAR(36) PRIMARY KEY,username VARCHAR(60), hashed_pass VARCHAR(60) NOT NULL, salt VARCHAR(50) NOT NULL);")

	if err != nil {
		panic(err)
		return
	}

	columns := []Column{
		NewVarcharColumn("username", "Test"),
	}

	_, err = repo.ExecuteQuery("SELECT uuid, username, hashed_pass, salt FROM UserTable WHERE username = ?", columns, QueryOptions{
		OnComplete: func(rows *sql.Rows) {
			defer rows.Close()

			for rows.Next() {

				var username string
				var uid uuid.UUID
				var hashedPassword string
				var salt string

				err := rows.Scan(&uid, &username, &hashedPassword, &salt)

				if err != nil {
					fmt.Printf("Scan error: %s\n", err)
					continue
				}

				fmt.Printf("uid: %s,\nusername: %s,\npassword: %s,\nsalt: %s\n", uid, username, hashedPassword, salt)

			}
			// Check for any error that occurred during iteration
			if err := rows.Err(); err != nil {
				fmt.Println("Rows iteration error:", err)
			}

		},
	})
	if err != nil {
		return
	}

}
