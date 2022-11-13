package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	User_id    int
	Name       string
	Email      string
	Password   string
	Deactive   int
	User_level string
}
type Posts struct {
	Post_id      int
	User_id      int
	Heading      string
	Body         string
	Closed_user  int
	Closed_admin int
	Closed_date  string
	Insert_time  string
	Update_time  string
	Image        string
}

type Comments struct {
	Comment_id  int
	Post_id     int
	User_id     int
	Body        string
	Insert_time string
}

type Categories struct {
	Category_id   int
	Category_Name string
	Closed        int
}

type Reaction struct {
	User_id    int
	Post_id    int
	Comment_id int
	Reaction   string
}

type PostCategory struct {
	Category_id int
	Post_id     int
}
type UserLevel struct {
	User_level    string
	value 	   int
}
func (u *Users) GetName() string {
	return u.Name
}

func (u *Users) GetEmail() string {
	return u.Email
}

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
	forumdb, err := sql.Open("sqlite3", "./"+databaseFile)

	if err != nil {
		return nil, err
	}

	return forumdb, nil
}

// function check if database exists, if not it creates it, if it does it opens it
func DatabaseExist() (*sql.DB, error) {
	newDb := false
	databaseFile := "forum-db.db"
	_, err := os.Stat(databaseFile)
	if os.IsNotExist(err) {
		fmt.Println("Creating the forum database ...")
		file, err := os.Create(databaseFile) // Create Sqlite file
		if err != nil {
			return nil, err
		}
		file.Close()
		fmt.Println("database created")

		newDb = true
	} else if err != nil {
		return nil, err
	}
	forumdb, err := sql.Open("sqlite3", "./"+databaseFile) // Open the created Sqlite3 File
	if err != nil {
		return nil, err
	}

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

func exampleDbData(forumdb *sql.DB) {
	insertUsers(forumdb, "peter", "peter@gritlab.ax", "bachelor", 0)
	insertUsers(forumdb, "aidran", "aidran@gritlab.ax", "younger", 1)
	insertUsers(forumdb, "tosin", "tosin@gritlab.ax", "kakkalla", 1)
	insertUsers(forumdb, "christian", "christain@gritlab.ax", "kingofhanko", 0)
	insertUsers(forumdb, "taneli", "tvntvn@gritlab.ax", "kakka", 1)
	insertUsers(forumdb, "jussi", "jussi@gritlab.ax", "kakka", 0)
	insertUsers(forumdb, "andre", "andre@gritlab.ax", "kakka923rfg", 0)
	insertPost(forumdb, 1, "football", "football is a great sport", "2017-01-01 00:00:00", "")
	insertPost(forumdb, 3, "Improving Data Science", "I wished to rebuild the whole data science world", "2017-01-01 00:00:45", "")
	insertPost(forumdb, 6, "Eating at a restaurants", "Me gusta is the best burger in Mariehamn", "2017-01-01 00:23:00", "")
	insertPost(forumdb, 1, "Hiking in the westlands", "Johannes is the best hicker in gritlab", "2017-01-01 00:60:00", "")
	insertCategory(forumdb, "sport")
	insertCategory(forumdb, "food")
	insertCategory(forumdb, "hiking")
	insertCategory(forumdb, "data science")
	insertCategory(forumdb, "programming")
	insertCategory(forumdb, "music")
	insertCategory(forumdb, "movies")
	insertCategory(forumdb, "books")
	insertComment(forumdb, 1, 1, "I agree with you", "2017-01-01 00:00:00")
	insertComment(forumdb, 2, 1, "I do not agree with you", "2017-01-01 00:00:00")
	insertReaction(forumdb, 1, 1, 1, "😀")
	insertReaction(forumdb, 1, 1, 2, "💩")
	insertPostCategory(forumdb, 1, 1)
}

func QueryResultDisplay(db *sql.DB) {
	row, err := db.Query(`
		SELECT users.name, posts.heading
		FROM users
		INNER JOIN posts ON users.user_id=posts.user_id
		WHERE users.user_id == 1;`)
	// Query the Database
	if err != nil {
		fmt.Println(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records
		var name string
		var heading string
		row.Scan(&name, &heading)                 // Fetch the record
		fmt.Println("user: ", name, "|", heading) // Print the record
	}
}
