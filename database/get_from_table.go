package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Get users from database
func GetUsers(db *sql.DB, userData map[string]string) ([]Users, error) {
	query := "select * from users WHERE"
	count := 0
	fmt.Println(userData)
	for k, v := range userData {
		if k == "password" {
			return nil, errors.New("password is not a valid search parameter")
		}
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'" + " COLLATE NOCASE"
		}
		count++
	}
	fmt.Println(query)
	var users []Users
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var user Users
		if err := rows.Scan(&user.User_id, &user.Name, &user.Email, &user.Password, &user.Deactive, &user.User_level); // Fetch the record
		err != nil {
			fmt.Println(err)
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Get posts from db
func GetPosts(db *sql.DB, postData map[string]string) ([]Posts, error) {
	query := "select * from posts WHERE"
	count := 0
	for k, v := range postData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var posts []Posts
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var post Posts
		if err := rows.Scan(&post.Post_id, &post.User_id, &post.Heading, &post.Body, &post.Closed_user, &post.Closed_admin, &post.Closed_date, &post.Insert_time, &post.Update_time, &post.Image); // Fetch the record
		err != nil {
			fmt.Println(err)
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// Get comments from db
func GetComments(db *sql.DB, commentData map[string]string) ([]Comments, error) {
	query := "select * from comments WHERE"
	count := 0
	for k, v := range commentData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var comments []Comments
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var comment Comments
		if err := rows.Scan(&comment.Comment_id, &comment.User_id, &comment.Post_id, &comment.Body, &comment.Insert_time); // Fetch the record
		err != nil {
			fmt.Println(err)
			return comments, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// get categories from db
func GetCategories(db *sql.DB, categoryData map[string]string) ([]Categories, error) {
	query := "select * from categories WHERE"
	count := 0
	for k, v := range categoryData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var categories []Categories
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var category Categories
		if err := rows.Scan(&category.Category_id, &category.Category_Name, &category.Closed); // Fetch the record
		err != nil {
			fmt.Println(err)
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// get reaction from db
func GetReaction(db *sql.DB, reactionData map[string]string) ([]Reaction, error) {
	query := "select * from reaction WHERE"
	count := 0
	for k, v := range reactionData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var reactions []Reaction
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var reaction Reaction
		if err := rows.Scan(&reaction.User_id, &reaction.Post_id, &reaction.Comment_id, &reaction.Reaction); // Fetch the record
		err != nil {
			fmt.Println(err)
			return reactions, err
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

// get post categories from db
func GetPostCategories(db *sql.DB, postCategoriesData map[string]string) ([]PostCategory, error) {
	query := "select * from postsCategory WHERE"
	count := 0
	for k, v := range postCategoriesData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	var postCategories []PostCategory
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var postCategory PostCategory
		if err := rows.Scan(&postCategory.Category_id, &postCategory.Post_id); // Fetch the record
		err != nil {
			fmt.Println(err)
			return postCategories, err
		}
		postCategories = append(postCategories, postCategory)
	}
	return postCategories, nil
}

// get user_level from db

func GetUserLevel(db *sql.DB, userLevelData map[string]string) ([]UserLevel, error) {
	query := "select * from userLevel WHERE"
	count := 0
	for k, v := range userLevelData {
		if count > 0 {
			query += " AND "
		}
		if k == "free_query" {
			query += " " + v
		} else {
			query += " " + k + "='" + v + "'"
		}
		count++
	}
	fmt.Println(query)
	var userLevels []UserLevel
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() { // Iterate and fetch the records
		var userLevel UserLevel
		if err := rows.Scan(&userLevel.User_level, &userLevel.value); // Fetch the record
		err != nil {
			fmt.Println(err)
			return userLevels, err
		}
		userLevels = append(userLevels, userLevel)
	}
	return userLevels, nil
}

// for test remove at the end
// func test() {
// 	user := make(map[string]string)
// 	post := make(map[string]string)
// 	comment := make(map[string]string)
// 	category := make(map[string]string)
// 	reaction := make(map[string]string)
// 	postCategory := make(map[string]string)
// 	userLevel := make(map[string]string)
// 	user["free_query"] = "user_id=2 AND name LIKE '%p%'"
// 	post["free_query"] = "user_id=1 OR heading LIKE '%t%'"
// 	comment["free_query"] = "user_id=1 OR body LIKE '%t%'"
// 	category["free_query"] = "category_id=2 OR category_name LIKE '%t%'"
// 	reaction["free_query"] = "user_id=1 OR post_id=1"
// 	postCategory["free_query"] = "post_id=1 OR category_id=2"
// 	userLevel["free_query"] = "user_level=1 OR value=1"
// 	db, err := DbConnect()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	users, err := GetUsers(db, user)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	posts, err := GetPosts(db, post)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	comments, err := getComments(db, comment)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	categories, err := getCategories(db, category)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	reactions, err := getReaction(db, reaction)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	postCategories, err := getPostCategories(db, postCategory)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	userLevels, err := getUserLevel(db, userLevel)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("users are ", users)
// 	fmt.Println("posts are", posts)
// 	fmt.Println("comments are", comments)
// 	fmt.Println("categories are", categories)
// 	fmt.Println("reactions are", reactions)
// 	fmt.Println("postCategories are", postCategories)
// 	fmt.Println("userLevels are", userLevels)
// }

// func init() {
// 	test()
// }
