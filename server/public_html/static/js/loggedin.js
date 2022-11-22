window.onload = setUser();

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
  

  let newPost = {
    postHeading: userPostHeading.value,
    postBody: userPost.value,
  };

  await fetch("/addPost", {
    method: "POST",
    body: JSON.stringify(newPost),
  })
    .then((response) => response.json())
    .then((json) => {
      console.log(json);
    });

  // create new post in DOM (old)
  const postDiv = document.createElement("div");
  postDiv.classList.add("border", "rounded", "content", "mx-auto", "col-8", "mb-2");
  postDiv.innerHTML = `<section class="row" id="post_section">
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
                      Create time ex (12:37)
                  </div>
                  <div class="col-6 order-1 text-end" id="post_mod_time">
                      Update time ex (12:53)
                  </div>
              </div>
          </div>
      </div>
  </div>

  <div class="offset-1 py-1">
      <div class="col-12 mb-2">
          <div class="row">
              <div class="mx-1" id="post_reactions">
                  <button class="bg-dark border rounded-start">👍
                      <span class="badge text-info">10</span>
                  </button>
                  <button class="bg-dark border rounded-end">👎
                      <span class="badge text-info">5</span>
                  </button>
                  <p class="text-info"># Comments</p>
              </div>
          </div>
      </div>
      <div class="col-10 justify-content-center mx-2 mb-2" id="user_comment">
        <div class="row">
          <div class="col-1>
            <img class="rounded-circle center-block" style="max-width: 55px" src="static/images/raccoon.jpeg" id="user_pic"></img>
          </div>
          <div class="col-10 text-start">
            <div class="input-group">
                <textarea
                    class="bg-dark border-info rounded text-light px-2 w-75"
                    style="resize:none;"
                    placeholder="Write a comment"></textarea>
                  <div class="input-group-append mx-2">
                    <button
                      class="btn bg-info text-dark mt-2"
                      type="button"
                      formaction="/addComment"
                      id="add_post_comment">
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
