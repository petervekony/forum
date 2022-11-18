window.onload = setUser();

async function setUser() {
  const userPic = document.getElementById("user_pic");
  const userName = document.getElementById("user_name");

  await fetch("/getUser", {
    method: "POST",
  })
    .then((response) => response.json())
    .then((json) => {
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

  // create new post in DOM
  const postSection = document.getElementById("post_section");

  const postDiv = document.createElement("div");
  // postDiv.id = postJSON.post_id;

  const postBody = document.createElement("div");
  const postBodyText = document.createElement("div");
  const postBodyTimeRow = document.createElement("div");
  const postHeading = document.createElement("div");
  const postInsertTime = document.createElement("div");
  const postModTime = document.createElement("div");
  const postImage = document.createElement("div");
  const postReactionsRow = document.createElement("div");
  const postReactions = document.createElement("div");
  const postLike = document.createElement("button");
  const postDislike = document.createElement("button");
  const postHeart = document.createElement("button");

  // const postLikeNum = document.createElement("p");

  postDiv.classList.add(
    "col-8",
    "offset-2",
    "text-white",
    "border",
    "rounded",
    "my-2",
    "py-2"
  );

  postBody.classList.add("col-10", "offset-1");
  postBodyText.classList.add("text-justify", "my-2");
  postBodyTimeRow.classList.add("row", "text-secondary");
  postReactionsRow.classList.add("row");
  postHeading.classList.add("col-10", "offset-1", "my-2");
  postInsertTime.classList.add("col-6", "order-0", "text-left");
  postModTime.classList.add("col-6", "order-1", "text-end");
  postReactions.classList.add("col-5", "mx-1");
  postImage.classList.add("border", "bg-info", "text-center");
  postLike.classList.add("btn", "border", "rounded", "bg-dark");
  postHeart.classList.add("btn", "border", "rounded", "bg-dark");
  postDislike.classList.add("btn", "border", "rounded", "bg-dark");

  // postLikeNum.classList.add("mx-1");

  postBodyText.textContent = `cookie monster just posted: ${userPost.value}`;
  // postHeading.textContent = postJSON.heading;
  // postInsertTime.textContent = postJSON.insert_time;
  // postModTime.textContent = postJSON.update_time;
  // postReactions.textContent = postJSON.post_reactions;
  // postImage.textContent = `<img src="/server/public_html/statis/images/${postJSON.image}">`;
  postLike.textContent = "👍";
  postDislike.textContent = "👎";
  postHeart.textContent = "❤️";
  // postLikeNum.textContent = "0";
  // test

  postSection.prepend(postDiv);
  postDiv.appendChild(postHeading);
  postDiv.appendChild(postBody);
  postBody.appendChild(postImage);
  postBody.appendChild(postBodyText);
  postBody.appendChild(postReactionsRow);
  postBodyText.appendChild(postBodyTimeRow);
  postBodyTimeRow.appendChild(postInsertTime);
  postBodyTimeRow.appendChild(postModTime);
  postReactionsRow.appendChild(postReactions);
  postReactions.appendChild(postLike);
  // postReactions.appendChild(postLikeNum);
  postReactions.appendChild(postDislike);
  postReactions.appendChild(postHeart);

  // let newPost = {
  //   sessionID: sessionID,
  //   uID: uID,
  //   postContent: userPost.value,
  //   postImage: userImage.value,
  // };

  // await fetch("/addPost", {
  //   method: "POST",
  //   body: JSON.stringify(newPost),
  // }).then((response) => response.json()).then((json) => {
  //   console.log(json)
  // });
  userPost.value = "";
}
