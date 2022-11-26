package server

import (
	"encoding/json"
	"gritface/database"
	"sort"
	"strconv"
)

type JSONData struct {
	Post_id      int                  `json:"post_id"`
	User_id      int                  `json:"user_id"`
	Heading      string               `json:"heading"`
	Body         string               `json:"body"`
	Closed_user  int                  `json:"closed_user"`
	Closed_admin int                  `json:"closed_admin"`
	Closed_date  string               `json:"closed_date"`
	Insert_time  string               `json:"insert_time"`
	Update_time  string               `json:"update_time"`
	Image        string               `json:"image"`
	Comments     map[int]JSONComments `json:"comments"`
	Categories	 []string							`json:"categories"`
	Reactions		 []map[int]string			`json:"reactions"`
	Username		 string								`json:"username"`
}

type JSONComments struct {
	CommentID int    							`json:"comment_id"`
	Post_id   int    							`json:"post_id"`
	User_id   int    							`json:"user_id"`
	Body      string 							`json:"body"`
	Reactions	[]map[int]string		`json:"reactions"`
	Username	string							`json:"username"`
}

func Retrieve20Posts() (string, error) {
	db, err := database.DbConnect()

	if err != nil {
		return "", err
	}

	structSlice := make(map[int]JSONData)
	query := "SELECT * FROM posts ORDER BY post_id DESC LIMIT 20"
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	nextQuery := ""
	for rows.Next() {
		rD := &JSONData{
			Comments: make(map[int]JSONComments),
		}
		err = rows.Scan(&rD.Post_id, &rD.User_id, &rD.Heading, &rD.Body, &rD.Closed_user, &rD.Closed_admin, &rD.Closed_date, &rD.Insert_time, &rD.Update_time, &rD.Image)
		if err != nil {
			return "", err
		}

		postId := &rD.Post_id

		// getting user's name
		currentUser := make(map[string]string)
		currentUser["user_id"] = strconv.Itoa(rD.User_id)
		users, err := database.GetUsers(db, currentUser)
		if err != nil {
			return "", err
		}
		rD.Username = users[0].Name;

		// getting post's categories
		currentPost := make(map[string]string)
		currentPost["post_id"] = strconv.Itoa(*postId)
		categories, err := database.GetPostCategories(db, currentPost)
		if err != nil {
			return "", err
		}
		var categoryNames []string
		for _, category := range categories {
			currentCategory := make(map[string]string)
			currentCategory["category_id"] = strconv.Itoa(category.Category_id)
			categoriesName, err := database.GetCategories(db, currentCategory)
			if err != nil {
				return "", err
			}
			categoryNames = append(categoryNames, categoriesName[0].Category_Name)
		}
		rD.Categories = categoryNames

		// getting post's reactions
		reactions, err := database.GetReaction(db, currentPost)
		if err != nil {
			return "", err
		}
		for _, reaction := range reactions {
			if reaction.Comment_id == 0 {
				userReaction := make(map[int]string)
				userReaction[reaction.User_id] = reaction.Reaction_id
				rD.Reactions = append(rD.Reactions, userReaction)
			}
		}

		structSlice[*postId] = *rD

		thisPostId := &rD.Post_id
		nextQuery += " OR post_id=" + strconv.Itoa(*thisPostId)
	}

	// Query comments
	query = "SELECT comment_id, post_id, user_id, body FROM comments WHERE " + nextQuery[4:]
	rows, err = db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		row := &JSONComments{}
		err = rows.Scan(&row.CommentID, &row.Post_id, &row.User_id, &row.Body)
		if err != nil {
			return "", err
		}
		currentUser := make(map[string]string)
		currentUser["user_id"] = strconv.Itoa(row.User_id)
		users, err := database.GetUsers(db, currentUser)
		if err != nil {
			return "", err
		}
		row.Username = users[0].Name;


			thisPostId := &row.Post_id
			thisCommentId := &row.CommentID
			// getting reactions
			currentComment := make(map[string]string)
			currentComment["comment_id"] = strconv.Itoa(*thisCommentId)
			reactions, err := database.GetReaction(db, currentComment)
			if err != nil {
				return "", err
			}
			for _, reaction := range reactions {
				userReaction := make(map[int]string)
				userReaction[reaction.User_id] = reaction.Reaction_id
				row.Reactions = append(row.Reactions, userReaction)
			}

			structSlice[*thisPostId].Comments[row.CommentID] = *row
		}
	}
	// The output needs to be in a descending order (by post_id), so we save it into a sorted []JSONData
	sSlice := make([]JSONData, 0, len(structSlice))
	for _, value := range structSlice {
		sSlice = append(sSlice, value)
	}
	res, err := json.Marshal(structSlice)
	if err != nil {
		return "", err
	}

	// fmt.Println(structSlice)
	// fmt.Println(string(res))
	return string(res), nil
}
