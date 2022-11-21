package server

import (
	"database/sql"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"

	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"

	d "gritface/database"

	"golang.org/x/image/draw"
)

// function to process image from the request
func processImage(file multipart.File) (image.Image, error) {
	var img image.Image
	var err error
	//get file type
	// to sniff the content type only the first
	// 512 bytes are used.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	// this will return a image/...
	filetype := http.DetectContentType(buffer)

	// reset the file pointer to the beginning
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	// we just slice the string to get the extension
	ext := filetype[6:]
	if ext != "jpeg" && ext != "png" && ext != "gif" && ext != "jpg" {
		return nil, errors.New("Invalid file type")
	}
	// check the file type for decoding
	switch ext {
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	case "jpeg", "jpg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
	case "gif":
		img, err = gif.Decode(file)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported image format")
	}
	return img, nil
}

// resize image reduces the size of the image
func resizeImage(img image.Image) error {
	imageHeight := 120
	divisor := img.Bounds().Max.Y / imageHeight
	width := img.Bounds().Max.X / divisor
	resized := image.NewRGBA(image.Rect(0, 0, width, imageHeight))
	// apporxbilinear gives a close quality to the original image, other methds trades between quality and speed
	draw.ApproxBiLinear.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)
	return nil
}

// Given image is 5120x2880
// Want to have it 2560x1440
func SaveImg(filePath string, img image.Image) (string, error) {
	imgFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	defer imgFile.Close()
	jpeg.Encode(imgFile, img, nil)
	return filePath, nil
}

// function to upload image
// we do not need any response writer just  request
// this is the one that will be called from the handler
func uploadImageHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != "POST" {
		http.Error(w, "bad Request Error", http.StatusBadRequest)
		return "", errors.New("method not allowed")
	}
	// max upload size is 20mb
	// this we can adjust depending on the size of the image we want
	err := r.ParseMultipartForm(20000 << 10)
	if err != nil {
		http.Error(w, "File size bigger than 20MB", http.StatusBadRequest)
		return "", errors.New("image size too bigs. Please use an image less than 20MB in size")
	}
	// get the file from the request
	file, fileHeader, err := r.FormFile("profileImage")
	if err != nil {
		fmt.Println("invalid size", err)
		return "", err
	}
	defer file.Close()
	img, err := processImage(file)
	if err != nil {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return "", err
	}
	err = resizeImage(img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	// destination path to store the image
	dstPath := "server/public_html/static/images/profile_imgs/" + fileHeader.Filename // need a better way to do store files
	// save the image to the path
	filePath, err := SaveImg(dstPath, img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	return filePath, nil
}

// store image path to db
func storeImageToDB(db *sql.DB, filePath string, userID int) error {
	// store the image path to the database
	_, err := db.Exec("UPDATE users SET profile_img = $1 WHERE id = $2", filePath, userID)
	if err != nil {
		return err
	}
	return nil
}

// function to change profile image
func changeProfilePicture(w http.ResponseWriter, r *http.Request) {
	// get user id from the request
	// this is the id of the user that is logged in
	// we will use this to get the user from the database
	// and update the profile image
	user_id := r.FormValue("user_id")
	// connect to db
	db, err := d.DbConnect()
	if err != nil {
		// do something
	}
	loginUser := make(map[string]string)
	loginUser["user_id"] = user_id
	// get user from the database
	user, err := d.GetUsers(db, loginUser)
	if err != nil {
		// do something
	}
	// if user is not found, nothing to change
	if len(user) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
	}
	// path of the original image
	fmt.Println(user[0].Profile_image)
	// get the image from the request
	// this is the image that the user wants to use as profile image
	// we will save this image to the server and update the user profile image
	// in the database
	imgName, err := uploadImageHandler(w, r)
	if err != nil {
		// do something
	}
	// update the user profile image
	user[0].Profile_image = imgName
	// update the user in the database
	err = d.UpdateUserData(db, loginUser, user_id)
	if err != nil {
		// do something
	}
}

/* func uploadProfileImage(r *http.Request) (string, error) {
	// max upload size is 1mb
	r.ParseMultipartForm(1000 << 10)
	// get the file from the form data
	file, header, err := r.FormFile("profileImage")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// create a file to save the image
	f, err := os.Create("server/public_html/static/images/profile_imgs" + header.Filename)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return header.Filename, nil
} */
