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
      userPic.innerHTML = `<img src="${json.Image}">`;
      userName.textContent = json.Username;
    });

  // for testing
  // userPic.textContent = "cookie monster's pic";
  // userName.textContent = "cookie monster";

  const addPostBtn = document.getElementById("add_post_button");
  addPostBtn.addEventListener("click", newPost);
}

async function setCategories() {
  await fetch("/getCategories", {
    method: "POST",
  })
    .then((response) => response.json())
    .then(function (json) {
      const catsList = document.getElementById("postCats");
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
    });
}

async function newPost() {
  const userPost = document.getElementById("user_post");
  if (!userPost.value) {
    console.log("Post is empty");
    return;
  }
  console.log(`New post button clicked and value is ${userPost.value}`);

  const userPostHeading = document.getElementById("user_post_title");
  if (!userPostHeading.value) {
    console.log("Post title is empty");
    return;
  }

  let categoriesSelected = [];
  let catInnerHTML = "";
  const postCats = document
    .getElementById("postCats")
    .querySelectorAll("input:checked");
  postCats.forEach((x) => {
    categoriesSelected.push(x.value);
    catInnerHTML += `#${x.value} `;
  });

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
      postID = json.message;
    });

  // create new post in DOM (old)
  const postDiv = document.createElement("div");
  postDiv.classList.add(
    "border",
    "rounded",
    "content",
    "mx-auto",
    "col-8",
    "mb-2"
  );
  postDiv.id = postID;
  const username = document.getElementById("user_name").textContent;
  postDiv.innerHTML = `<section class="row" id="post_section">
  <p class="text-start mx-2 text-info">${username} ${catInnerHTML}</p>
  <div data-bs-target="#collapse_post_comments" data-bs-toggle="collapse">
      <div class="text-white rounded my-2 py-2" id="post_div">
          <div class="col-11 offset-1 my-1" id="post_heading">
              ${userPostHeading.value}
          </div>
          <div class="col-10 offset-1" id="post_body">
              <div class="border bg-info text-center" id="post_image">Testing image"</div>
              <div class="text-justify my-2">
                  ${userPost.value}
              </div>
              <div class="row text-secondary">
                  <div class="col-6 order-0 text-left" id="post_insert_time">
                      Created just now..
                  </div>
                  <div class="col-6 order-1 text-end" id="post_mod_time">
                  </div>
              </div>
          </div>
      </div>
  </div>

  <div class="offset-1 py-1">
      <div class="col-12 mb-2">
          <div class="row">
              <div class="mx-1" id="post_reactions">
                  <button class="bg-dark border rounded-start">⬆️<span
                          class="badge text-info">0</span></button>
                  <button class="bg-dark border rounded-end">⬇️<span class="badge text-info">0</span></button>
                  <p class="text-info">0 Comments</p>
              </div>
          </div>
          <div class="text-secondary">
          <p>{$categories} #IRL #Shitpost #depression</p>
        </div>
      </div>
      <div class="col-10 justify-content-center mx-3 mb-2" id="user_comment">
      <div class="row">
          <div class="col-1 mx-2">
                <img class="rounded-circle center-block" style="max-width: 55px; border: 2px solid #54B4D3;" src="static/images/raccoon.jpeg" id="user_pic">
          </div>
          <div class="col-10 text-start">
              <div class="input-group">
                  <textarea
                      class="bg-dark border-info rounded text-light px-2 w-75"
                      class="form-control"
                      style="resize:none;"
                      id="newComment"
                      placeholder="Write a comment"></textarea>
                  <div class="input-group-append mx-2">
                    <button
                      class="btn bg-info text-dark mt-2"
                      type="button"
                      onclick="addComment(${postDiv.id})">
                      Comment
                    </button>
                  </div>
                </div>
          </div>
      </div>
  </div>
</section>`;
  const container = document.getElementById("container");
  container.prepend(postDiv);
  userPost.value = "";
  userPostHeading.value = "";
}

async function addComment(id) {
  const postDiv = document.getElementById(id);
  const newComment = postDiv.querySelector("#newComment");
  if (!newComment.value) {
    console.log("Comment is empty");
    return;
  }
  console.log(`New comment button clicked and value is ${newComment.value}`);
  let comment = {
    postComment: newComment.value,
    postID: id,
  };
  console.log(id);
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
  commentDiv.classList.add("row", "my-3", "ms-auto");
  commentDiv.id = commentID;
  const userPic = document.getElementById("user_pic");
  const userName = document.getElementById("user_name");
  commentDiv.innerHTML = `<div class="col-1 mx-2">
      <img class="rounded-circle" style="max-width: 120%; border: 2px solid #54B4D3;" src="${userPic.getAttribute(
        "src"
      )}">
    </div>
    <div class="col-8 border rounded bg-secondary" id="post_comments">
    <p class="text-info pt-2">${userName.textContent}</p>
      ${newComment.value}
      <div class="row">
      <div class="text-end" id="comment_reactions">
        <button class="btn btn-dark rounded-start">⬆️
            <span class="badge text-info">0</span>
        </button>
        <button class="btn btn-dark rounded-end">⬇️
            <span class="badge text-info">0</span>
        </button>
      </div>
    </div>
      </div>
    </div>`;
  const commentsDiv = postDiv.querySelector(`#collapse_post_comments${id}`);
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
