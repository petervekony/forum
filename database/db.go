package database

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	logger "gritface/log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// hash password returned the password string as a hash to be stored in the database
// this is doen for security reasons
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// check passwords checks if the password provided by the user matches the one in the database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Connect to database
func DbConnect() (*sql.DB, error) {
	databaseFile := "forum-db.db"
	forumdb, err := sql.Open("sqlite3", "./"+databaseFile+"?_auth&_auth_user=forum&_auth_pass=forum&_auth_crypt=sha1")

	if err != nil {
		return nil, err
	}

	// Enable foreign key contraints
	enableContraints := `PRAGMA foreign_keys = ON;`

	enableContraintsQuery, err := forumdb.Prepare(enableContraints) // Prepare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	enableContraintsQuery.Exec()

	return forumdb, nil
}

// function check if database exists, if not it creates it, if it does it opens it
func DatabaseExist() (*sql.DB, error) {
	newDb := false
	databaseFile := "forum-db.s3db"
	_, err := os.Stat(databaseFile)
	if os.IsNotExist(err) {
		fmt.Println("Creating the forum database ...")
		file, err := os.Create(databaseFile) // Create Sqlite file
		if err != nil {
			return nil, err
		}
		file.Close()

		logger.WTL("Database created", true)
		newDb = true
	} else if err != nil {
		return nil, err
	}
	forumdb, err := sql.Open("sqlite3", "./"+databaseFile)
	// Open the created Sqlite3 File
	if err != nil {
		logger.WTL("Database could not be opened", false)
		return nil, err
	}
	conn, err := forumdb.Conn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if newDb {
		err = createTable(forumdb) // Create Database Tables
		if err != nil {
			fmt.Println(err.Error())
		} else {
			// INSERT RECORDS
			exampleDbData(forumdb)
		}
	}

	var requiredTables = map[string]bool{"users": false, "posts": false, "comments": false, "categories": false}
	tables, table_check := forumdb.Query("select name from sqlite_master where type='table' and name not like 'sqlite_%'")

	if table_check == nil {
		for tables.Next() { // Iterate and fetch the records
			var name string
			tables.Scan(&name)
			requiredTables[name] = true // Fetch the record
		}
		for _, value := range requiredTables {
			if !value {
				forumdb.Close() // Close connection to current database
				reader := bufio.NewReader(os.Stdin)
			handleInvalidDatabase:
				fmt.Println("Existing database is not working with this version of forum.\nWould you like to delete current database (y) or rename current database (r), quit (q)?:")
				uinput, _ := reader.ReadString('\n')
				uinput = strings.Trim(uinput, "\n")

				if uinput == "y" {
					os.Remove(databaseFile)
				} else if uinput == "r" {
				renameCurrentDB:
					fmt.Print("Rename to: ")
					uinput, _ = reader.ReadString('\n')
					uinput = strings.Trim(uinput, "\n")
					if len(uinput) < 1 {
						fmt.Println("Name can not be empty")
						goto renameCurrentDB
					} else if len(uinput) < 4 { // Add .db to the end as a string 3 char long can not hold name + .db
						uinput += ".db"
					} else if uinput[len(uinput)-4:] != ".db" {
						uinput += ".db"
					}
					os.Rename(databaseFile, uinput)
					fmt.Println("Database renamed to " + uinput)
				} else if uinput == "q" {
					fmt.Println("Exiting....")
					os.Exit(0)
				} else {
					goto handleInvalidDatabase
				}
				return DatabaseExist()
			}
		}
	} else {
		fmt.Println("table not there")
	}
	return forumdb, nil
}

// remove this when cleaning up
func exampleDbData(forumdb *sql.DB) {
	InsertUsers(forumdb, "peter", "'; DROP TABLE users;'", "bachelor", 0)
	InsertUsers(forumdb, "aidran", "aidran@gritlab.ax", "younger", 1)
	InsertUsers(forumdb, "tosin", "tosin@gritlab.ax", "kakkalla", 1)
	InsertUsers(forumdb, "christian", "christain@gritlab.ax", "kingofhanko", 0)
	InsertUsers(forumdb, "taneli", "tvntvn@gritlab.ax", "kakka", 1)
	InsertUsers(forumdb, "jussi", "jussi@gritlab.ax", "kakka", 0)
	InsertUsers(forumdb, "andre", "andre@gritlab.ax", "kakka923rfg", 0)
	InsertPost(forumdb, 1, "football", "football is a great sport", "2017-01-01 00:00:00", "")
	InsertPost(forumdb, 3, "Improving Data Science", "I wished to rebuild the whole data science world", "2017-01-01 00:00:45", "")
	InsertPost(forumdb, 6, "Eating at a restaurants", "Me gusta is the best burger in Mariehamn", "2017-01-01 00:23:00", "")
	InsertPost(forumdb, 1, "Hiking in the westlands", "Johannes is the best hicker in gritlab", "2017-01-01 00:60:00", "")
	InsertCategory(forumdb, "sport")
	InsertCategory(forumdb, "food")
	InsertCategory(forumdb, "hiking")
	InsertCategory(forumdb, "data science")
	InsertCategory(forumdb, "programming")
	InsertCategory(forumdb, "music")
	InsertCategory(forumdb, "movies")
	InsertCategory(forumdb, "books")
	InsertComment(forumdb, 1, 1, "I agree with you", "2017-01-01 00:00:00")
	InsertComment(forumdb, 2, 1, "I do not agree with you", "2017-01-01 00:00:00")
	InsertReaction(forumdb, 1, 1, 1, "1")
	InsertReaction(forumdb, 1, 2, 2, "1")
	InsertReaction(forumdb, 8, 2, 2, "1")
	InsertReaction(forumdb, 1, 1, 3, "2")
	InsertReaction(forumdb, 1, 1, 4, "2")
	InsertReaction(forumdb, 1, 1, 5, "2")
	InsertReaction(forumdb, 1, 1, 0, "2")
	InsertReaction(forumdb, 1, 2, 0, "2")
	InsertPostCategory(forumdb, 1, 1)
}
