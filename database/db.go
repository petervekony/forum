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

func createTable(db *sql.DB) error {
	// users table
	createUsersTable := `CREATE TABLE users (
		"user_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,		
		"name" TEXT NOT NULL UNIQUE,
		"email" TEXT NOT NULL UNIQUE,
		"password" TEXT NOT NULL,
		"deactive" INTEGER DEFAULT 0,
		"user_level" INTEGER
	  );`

	usersStatement, err := db.Prepare(createUsersTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	usersStatement.Exec() // Execute SQL Statements

	// posts tables
	createPostsTable := `CREATE TABLE posts (
		"post_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
		"user_id" INTEGER NOT NULL,		
		"heading" TEXT NOT NULL,
		"body" TEXT NOT NULL,
		"closed_user" INTEGER default '0',
		"closed_admin" INTEGER default '0',
		"closed_date" TEXT DEFAULT '',
		"insert_time" TEXT NOT NULL,
		"update_time" TEXT NOT NULL DEFAULT '', 
		"image" TEXT
	  );`
	postsStatement, err := db.Prepare(createPostsTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	postsStatement.Exec() // Execute SQL Statements

	// comments table
	createcommentsTable := `CREATE TABLE comments (
		"comment_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
		"post_id" INTEGER NOT NULL,
		"user_id" INTEGER NOT NULL,
		"body" TEXT NOT NULL,
		"insert_time" TEXT NOT NULL
	  );`

	commentsStatement, err := db.Prepare(createcommentsTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	commentsStatement.Exec() // Execute SQL Statements

	// categories table
	createcategoriesTable := `CREATE TABLE categories (
		"category_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
		"category_name" TEXT NOT NULL UNIQUE,
		"closed" INTEGER default 0
	  );`
	categoriesStatement, err := db.Prepare(createcategoriesTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	categoriesStatement.Exec() // Execute SQL Statements

	// reaction table
	createreactionTable := `CREATE TABLE reaction (
		"user_id" INTEGER NOT NULL,
		"post_id" INTEGER NOT NULL,
		"comment_id" INTEGER NOT NULL,
		"reaction" text NOT NULL,
		PRIMARY KEY (user_id, post_id, comment_id)
	  );`
	reactionStatement, err := db.Prepare(createreactionTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	reactionStatement.Exec() // Execute SQL Statements

	//post category table
	createpostscategoryTable := `CREATE TABLE postscategory (
		"category_id" INTEGER NOT NULL,
		"post_id" INTEGER NOT NULL
	  );`
	postcategoryStatement, err := db.Prepare(createpostscategoryTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	postcategoryStatement.Exec() // Execute SQL Statements

	createUserLevelTable := `CREATE TABLE user_level (
		"user_level" TEXT NOT NULL UNIQUE,
		VALUE INTEGER NOT NULL,
		PRIMARY KEY (user_level, value)
	)`
	userlevelStatement, err := db.Prepare(createUserLevelTable) // Prepare SQL Statement
	if err != nil {
		return err
	}
	userlevelStatement.Exec() // Execute SQL Statements

	return nil
}

// insertUsers function inserts a record in the users table
func insertUsers(db *sql.DB, name string, email string, password string, user_level int) {
	Password, _ := HashPassword(password)
	insertUsers := `INSERT INTO users(name, email, Password, user_level) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertUsers) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(name, email, Password, user_level) // Execute statement with parameters
	if err != nil {
		fmt.Println(err.Error())
	}
}

// function adds posts to the database
func insertPost(db *sql.DB, user_id int, heading string, body string, insert_time string, image string) {
	insertPost := `INSERT INTO posts(user_id, heading, body, insert_time, image) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertPost) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(user_id, heading, body, insert_time, image) // Execute statement with parameters
	if err != nil {
		fmt.Println(err.Error())
	}
}

// function inserts categories into the database
func insertCategory(db *sql.DB, category_name string) {
	insertCategory := `INSERT INTO categories(category_name) VALUES (?)`
	statement, err := db.Prepare(insertCategory) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(category_name) // Execute statement with parameters
	if err != nil {
		fmt.Println(err.Error())
	}
}

// function inserts comments into the database
func insertComment(db *sql.DB, post_id int, user_id int, body string, insert_time string) {
	insertComment := `INSERT INTO comments(post_id, user_id, body, insert_time) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertComment) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(post_id, user_id, body, insert_time) // Execute statement with parameters
	if err != nil {
		fmt.Println(err.Error())
	}
}

// function inserts reaction into the database
func insertReaction(db *sql.DB, user_id int, post_id int, comment_id int, reaction string) {
	insertReaction := `INSERT INTO reaction(user_id, post_id, comment_id, reaction) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertReaction) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(user_id, post_id, comment_id, reaction) // Execute statement with parameters
	if err != nil {
		fmt.Println(err.Error())
	}
}

// function inserts post category into the database
func insertPostCategory(db *sql.DB, post_id int, category_id int) {
	insertPostCategory := `INSERT INTO postscategory(post_id, category_id) VALUES (?, ?)`
	statement, err := db.Prepare(insertPostCategory) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = statement.Exec(post_id, category_id) // Execute statement with parameters
	if err != nil {
		fmt.Println(err.Error())
	}
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

// function check if database exists, if not it creates it, if it does it opens it
func DatabaseExist() (*sql.DB, error) {
	newDb := false
	databaseFile := "forum-db.db"
	if _, err := os.Stat(databaseFile); os.IsNotExist(err) {
		fmt.Println("Creating the forum database ...")
		file, err := os.Create(databaseFile) // Create Sqlite file
		if err != nil {
			return nil, err
		}
		file.Close()
		fmt.Println("database created")

		newDb = true
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
