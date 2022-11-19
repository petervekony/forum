package server

import (
	"errors"
	"fmt"
	"image"
	"mime/multipart"

	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"

	"golang.org/x/image/draw"
	d "gritface/database"
)

const size = 1000

// const maxSize = 20000 * 1024
const maxFileSize = size << 10

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
	// we just slice the string to get the extension
	ext := filetype[6:]
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

	// destination path to store the image
	dstPath := "server/public_html/static/images/posts_imgs/"
	// save the image to the path
	SaveImg(dstPath, resized)
	return nil
}

// Given image is 5120x2880
// Want to have it 2560x1440
func SaveImg(filePath string, img image.Image) {
	imgFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	defer imgFile.Close()
	jpeg.Encode(imgFile, img, nil)
}

// function to upload image
// we do not need any response writer just  request
// this is the one that will be called from the handler
func ImageUpload(r *http.Request) (string, error) {
	// max upload size is 20mb
	// this we can adjust depending on the size of the image we want
	r.ParseMultipartForm(20000 << 10)
	// get the file from the request
	file, _, err := r.FormFile("profileImage")
	if err != nil {
		fmt.Println("invalid size", err)
		return "", err
	}
	defer file.Close()
	img, err := processImage(file)
	if err != nil {
		return "", err
	}
	err = resizeImage(img)
	if err != nil {
		return "", err
	}
	return "Image successfully uploaded", nil
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
	fmt.Println(users.I)
	// get the image from the request
	// this is the image that the user wants to use as profile image
	// we will save this image to the server and update the user profile image
	// in the database
	_, err = ImageUpload(r)
	if err != nil {
		// do something
	}
	




	// call the image upload function
	_, err := ImageUpload(r)
	if err != nil {
		fmt.Println("error", err)
		return
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

/*
func Post(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxFileSize {
		if flusher, ok := w.(http.Flusher); ok {
			response := []byte("Request too large")
			w.Header().Set("Connection", "close")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(response)))
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write(response)
			flusher.Flush()
		}
		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	err := r.ParseMultipartForm(1024)
	if err != nil {
		w.Write([]byte("File too large"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}

	dst, err := os.Create("./public_html/static/images/" + header.Filename)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, io.LimitReader(file, maxFileSize))
	if err != nil {
		panic(err)
	}

	if written == maxFileSize {
		w.Write([]byte("File too large"))
		return
	}
	w.Write([]byte("Success..."))
} */
