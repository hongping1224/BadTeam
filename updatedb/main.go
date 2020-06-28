package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {

	db, err := sql.Open("mysql", "admin:NanaDatabasePassword@tcp(badteam.ccz3kc9rn8lq.ap-southeast-1.rds.amazonaws.com:3306)/badteam")
	if err != nil {
		return fmt.Sprintf(" sql.Open Error: %v", err), nil
	}
	defer db.Close()

	if name.Name == "create" {
		err = CreateTable(db, "test", "id INT NOT NULL AUTO_INCREMENT,PRIMARY KEY(id), data INT, quo VARCHAR(30)")
		if err != nil {
			return fmt.Sprintf(" Create Table Error: %v", err), nil
		}
		return "Create Success", nil
	}

	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func CreateTable(db *sql.DB, tablename, colume string) error {
	q := fmt.Sprintf("CREATE TABLE %s(%s);", tablename, colume)
	rows, err := db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

/*
	endpoints := "badteam.ccz3kc9rn8lq.ap-southeast-1.rds.amazonaws.com"
	region := "ap-southeast-1a"
	dbuser := "admin"
	dbName := "badteam"
	sess := session.Must(session.NewSession())
	creds := stscreds.NewCredentials(sess, "lambda-vpc-role")
	authToken, err := rdsutils.BuildAuthToken(endpoints, region, dbuser, creds)
	if err != nil {
		return fmt.Sprintf("BuildAuthToken Error: %v", err), nil
	}
	dnsStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true", dbuser, authToken, endpoints, dbName)
	db, err := sql.Open("mysql", dnsStr)
	if err != nil {
		return fmt.Sprintf(" sql.Open Error: %v", err), nil
	}*/
//db.Stats()
