    # Welcome to GritFace Project
    This forum app is our solution of 01Edu's forum project. It consists of a message board where users can make posts and comments in multiple topics/categories, and leave likes and dislikes on each other's or their own posts and comments. Registered users can filter the posts according to four criteria:
    - posts created by the user
    - posts liked by the user
    - posts disliked by the user
    - categories

    # Usage
    To run the server, after cloning the repo to your local system, you can create an executable in the terminal from the main directory and run it using the following two commands:
    ```
    go build -o gritface main.go
    ./gritface
    ```
    If you want to run it without creating the executable, you can run the program directly:
    ```
    go run main.go
    ```
    After that, when you see the *Server listening on ':443'* output on your terminal, you can access the site on the *https://localhost* adress.

    ## Technologies
    The backend of the app is written in Go, using only the allowed packages, which consist of the standard library and the following:
    - [go-sqlite3](https://github.com/mattn/go-sqlite3)
    - [bcrypt](https://golang.org/x/crypto/bcrypt)
    - [UUID](https://github.com/gofrs/uuid)

    The frontend is written using Javascript, CSS, HTML and [Bootstrap](https://getbootstrap.com/).
    ### Docker
    If you want to run the app in a Docker container, you can use the Dockerfile in the root directory. 
    There is a bash script in the root directory, which builds the Docker image and runs it in a container. You can execute it using the following command:
    ```
    ./docker.sh
    ```
    To create the image manually, you can run the following:
    ```
    docker image build -f Dockerfile -t forum .
    ```
    After that, to run the container:
    ```
    docker container run -p 443:443 --detach --name gritface forum
    ```

    ### Javascript
    ### HTML
    ### CSS

    ## Project Layout/structure

    ## Website and Project Features

    // this should be in a list
    ### Signup
    ### Login
    ### Cookie based Session
    ### One active session per user
    ### Posts
    ### Comments
    ### Posts Like and Dislikes
    ### Comment Likes and Dislikes

    ## Third Party Packages

    // this shoudl be in a list





    ### File structure
    ### Database structure
    ### Server structure
    ### Frontend structure

    All your files and folders are presented as a tree in the file explorer. You can switch from one to another by clicking a file in the tree.

    ## 

    You can rename the current file by clicking the file name in the navigation bar or by clicking the **Rename** button in the file explorer.

    ## Contributors
    Adrain L, Christian D, Peter V, Tosin I, Taneli M