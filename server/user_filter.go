package server

import (
	"encoding/json"
	"fmt"
	d "gritface/database"
	"net/http"
	"sort"
	"strconv"
)

func userFilter(w http.ResponseWriter, r *http.Request, filter string, uid string) (string, error) {

	db, err := d.DbConnect()

	if err != nil {
		return "", err
	}
	filterValue := r.URL.Query()["filter"]
	useFilter := ""
	if len(filterValue) != 0 {
		useFilter = filterValue[0]
	}
	query := "SELECT posts.* from posts"
	switch useFilter {
	case "userPosts":
		query += " WHERE user_id=" + uid
	case "liked":
		// TODO: change "1" to uid
		query += " INNER JOIN reaction ON posts.post_id=reaction.post_id WHERE reaction.reaction='⬆️' AND reaction.user_id=" + "1" + " GROUP BY posts.post_id"
	case "category":
		cat := r.URL.Query()["cat"][0]
		fmt.Println("WHOLE FUCKING QUERY: ", r.URL.Query())
		query += " INNER JOIN postsCategory ON posts.post_id=postsCategory.post_id WHERE postsCategory.category_id=" + cat
		fmt.Println(query)
	default:
		data := &JSONData{
			Post_id:  0,
			User_id:  0,
			Heading:  "ERROR",
			Body:     "Invalid filter",
			Comments: make(map[int]JSONComments),
		}
		dataSlice := []JSONData{*data}
		dummy, _ := json.Marshal(dataSlice)
		return string(dummy), nil
	}
	structSlice := make(map[int]JSONData)
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	fmt.Println("Query: ", query)
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
		users, err := d.GetUsers(db, currentUser)
		if err != nil {
			return "", err
		}
		rD.Username = users[0].Name

		// getting post's categories
		currentPost := make(map[string]string)
		currentPost["post_id"] = strconv.Itoa(*postId)
		categories, err := d.GetPostCategories(db, currentPost)
		if err != nil {
			return "", err
		}
		var categoryNames []string
		for _, category := range categories {
			currentCategory := make(map[string]string)
			currentCategory["category_id"] = strconv.Itoa(category.Category_id)
			categoriesName, err := d.GetCategories(db, currentCategory)
			if err != nil {
				return "", err
			}
			categoryNames = append(categoryNames, categoriesName[0].Category_Name)
		}
		rD.Categories = categoryNames

		// getting post's reactions
		reactions, err := d.GetReaction(db, currentPost)
		if err != nil {
			return "", err
		}
		for _, reaction := range reactions {
			if reaction.Comment_id == 0 {
				userReaction := make(map[int]string)
				userReaction[reaction.User_id] = reaction.Reaction
				rD.Reactions = append(rD.Reactions, userReaction)
			}
		}

		structSlice[*postId] = *rD

		thisPostId := &rD.Post_id
		nextQuery += " OR post_id=" + strconv.Itoa(*thisPostId)
	}
	// Query comments
	if len(nextQuery) > 4 {
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
			users, err := d.GetUsers(db, currentUser)
			if err != nil {
				return "", err
			}
			row.Username = users[0].Name

			thisPostId := &row.Post_id
			thisCommentId := &row.CommentID
			// getting reactions
			currentComment := make(map[string]string)
			currentComment["comment_id"] = strconv.Itoa(*thisCommentId)
			reactions, err := d.GetReaction(db, currentComment)
			if err != nil {
				return "", err
			}
			for _, reaction := range reactions {
				userReaction := make(map[int]string)
				userReaction[reaction.User_id] = reaction.Reaction
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
	sort.Slice(sSlice, func(i, j int) bool { return sSlice[i].Post_id > sSlice[j].Post_id })

	res, err := json.Marshal(sSlice)
	if err != nil {
		return "", err
	}

	// fmt.Println(structSlice)
	return string(res), nil
}
