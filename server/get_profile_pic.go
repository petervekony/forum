package server

import (
	"errors"
	d "gritface/database"
)

func GetProfilePic(uid string) (string, error) {
	userMap := map[string]string{
		"user_id": uid,
	}
	forumdb, err := d.DbConnect()
	if err != nil {
		return "", err
	}
	user, err := d.GetUsers(forumdb, userMap)
	if err != nil {
		return "", err
	}
	if len(user) < 1 {
		return "", errors.New("ERROR: User not found")
	}
	return user[0].Profile_image, nil
}
