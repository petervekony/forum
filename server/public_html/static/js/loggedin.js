window.onload = setUser();
window.onload = setCategories();

async function setUser() {
  const userPic = document.getElementById("user_pic");
  const userName = document.getElementById("user_name");

  await fetch("/getUser", {
    method: "POST",
  })
    .then((response) => response.json())
    .then((json) => {
      if (!json.Username) {
        window.location.replace("/");
      }
      console.log(json);
      // userPic.innerHTML = `<img src="${json.Image}">`;
      userName.textContent = json.Username;
      userName.profile_image = json.profile_image;
    });

  const addPostBtn = document.getElementById("add_post_button");
  addPostBtn.addEventListener("click", newPost);
}

async function setCategories() {
  await fetch("/getCategories", {
    method: "POST",
  })
    .then((response) => response.json())
    .then(function (json) {
      console.log(json);
      const catsList = document.getElementById("postCats");
      const filterCatsList = document.getElementById("filterCats");
      const map = new Map(Object.entries(json));
      for (const [key, value] of map) {
        const catsItem = document.createElement("li");
        catsItem.innerHTML = `<a class="dropdown-item" href="#">
        <div class="form-check">
            <input class="form-check-input" type="checkbox" value="" id="${key}" />
            <label class="form-check-label" for="${key}">${value}</label>
        </div>
    </a>`;
        catsList.append(catsItem);
        const filterCatsItem = document.createElement("li");
        filterCatsItem.innerHTML = `<a class="dropdown-item bg-dark" href="#" onclick="initPage('/filtered?filter=category&cat=${key}')">
        <div class="col-12">
            <button class="col-12 btn btn-dark text-info border-info" value="${value}" id="fcheck${key}">${value}</button>
        </div>
    </a>`;
        filterCatsList.append(filterCatsItem);
      }
      /*
      json.map(function (category) {
        const catsItem = document.createElement("li");
        catsItem.innerHTML = `<a class="dropdown-item" href="#">
        <div class="form-check">
            <input class="form-check-input" type="checkbox" value="${category}" id="check1" />
            <label class="form-check-label" for="check1">${category}</label>
        </div>
    </a>`;
        catsList.append(catsItem);
      });
      const filterCatsList = document.getElementById("filterCats");
      json.map(function (category) {
        const catsItem = document.createElement("li");
        catsItem.innerHTML = `<a class="dropdown-item" href="#">
        <div class="form-check">
            <input class="btn btn-dark" type="button" value="${category}" id="check2" />
            <label class="form-check-label" for="check2">${category}</label>
        </div>
    </a>`;
        filterCatsList.append(catsItem);
      });
      */
    });
}

function isPostEmpty(userPost, userPostHeading) {
  if (!userPost.value) {
    console.log("Post is empty");
    return false;
  }
  if (!userPostHeading.value) {
    console.log("Post title is empty");
    return false;
  }
  return true
}

async function newPost() {
  const userPost = document.getElementById("user_post");
  const userPostHeading = document.getElementById("user_post_title");

  // validate empty post
  if (isPostEmpty(userPost, userPostHeading)) {
    console.log(`New post button clicked and value is ${userPost.value}`)
  } else {
    alert("Please enter all fields!");
    return
  }

  // get selected categories and form the text for in the front end in the new post
  let categoriesSelected = [];
  let catInnerHTML = "";
  const postCats = document
    .getElementById("postCats")
    .querySelectorAll("input:checked");
  postCats.forEach((x) => {
    categoriesSelected.push(x.value);
    catInnerHTML += `#${x.value} `;
  });

  // creating JSON and making request to the server for registering new post
  let newPost = {
    postHeading: userPostHeading.value,
    postBody: userPost.value,
    postCats: categoriesSelected,
  };

  let postID;
  await fetch("/addPost", {
    method: "POST",
    body: JSON.stringify(newPost),
  })
    .then((response) => response.json())
    .then((json) => {
      console.log(json);
      if (!json.status) {
        alert("No spamming posts please!")
        return
      }
      postID = json.message;
      initPage('/filtered?filter=postId&id=' + postID);
      userPost.value = "";
      userPostHeading.value = ""; 
    });
}

async function addComment(postID) {
  const postDiv = document.getElementById(postID);
  const newComment = postDiv.querySelector("#newComment");
  if (!newComment.value) {
    console.log("Comment is empty");
    return;
  }
  console.log(`New comment button clicked and value is ${newComment.value}`);
  let comment = {
    postComment: newComment.value,
    postID: postID,
  };
  console.log(postID);
  let commentID = 1;
  await fetch("/addComment", {
    method: "POST",
    body: JSON.stringify(comment),
  })
    .then((response) => response.json())
    .then((json) => {
      console.log(json);
      commentID = json.message;
    });
  // create new comment in DOM (old)
  const commentDiv = document.createElement("div");
  // commentDiv.classList.add("row", "mx-auto");
  commentDiv.postID = commentID;
  const userPic = document.getElementById("user_pic");
  const userName = document.getElementById("user_name");
  commentDiv.innerHTML = createCommentDiv(postID, commentID, userPic.getAttribute("src"), userName.textContent, newComment.value, 0, 0, 0, "Commented just now...");

  const commentsDiv = postDiv.querySelector(`#collapse_post_comments${postID}`);
  commentsDiv.classList.add("show");
  if (!commentsDiv) {
    console.log("broke down");
  } else {
    commentsDiv.prepend(commentDiv);
  }
  newComment.value = "";
  const number_of_comments = postDiv.querySelector("#number_of_comments");
  console.log(number_of_comments);
  number_of_comments.textContent =
    parseInt(number_of_comments.textContent) + 1 + " Comments";
}
