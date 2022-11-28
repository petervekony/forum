window.onload = initPage();

async function initPage(request = "/posts") {
  console.log("hello from initpage, " + request);
  await fetch(request)
    .then((response) => response.json())
    .then(function (json) {
      console.log(json)
      document.getElementById("container").innerHTML = "";
      createPosts(json);
    });
  if (request != "/posts") {
    document.getElementById("load_more_btn").style.display = "none";
  } else {
    document.getElementById("load_more_btn").style.display = "";
  }
}

function reactionButton(postId, commentId, reactId, reactCount, setActive=false) {
  let returndata = "";
  var reactIcon = `⬇️`;
  if(reactId == 1) {
    reactIcon = `⬆️`;
  }

  let addClass = "";
  if(setActive) {
    addClass += "active"
  }

  var setStyle = ""
  if(commentId > 0) {
    setStyle += "height:60%;";
    addClass += "px-0 py-0";
  } else {
    addClass += " border";
  }

  returndata += `<button class="btn btn-dark ${addClass}" style="${setStyle}" id="rbc${postId}_${commentId}_${reactId}" onclick="addReaction(${postId}, ${commentId}, ${reactId}, this)">${reactIcon}
                  <span class="badge text-info" id="rb${postId}_${commentId}_${reactId}">${reactCount}</span>
                  </button>`
  return returndata;
}

async function addReaction(postID, commentID, reactionID, targetButton) {
  await fetch(
    "/add_reaction?post_id=" +
      postID +
      "&comment_id=" +
      commentID +
      "&reaction_id=" +
      reactionID
  )
    .then((response) => response.json())
    .then(function (json) {
      console.log(json)
      for(let i=1; i < 3; i++) {
        document.getElementById("rb" + postID + "_" + commentID + "_" + i).innerHTML = json['rb'+i];
        document.getElementById("rbc" + postID + "_" + commentID + "_" + i).classList.remove("active")
      }
      if(json.userReaction > 0) {
        targetButton.classList.add("active");
      }
    });
}

async function signup() {
  const username = document.getElementById("signup_name").value;
  const email = document.getElementById("signup_email").value;
  const password = document.getElementById("signup_pass").value;
  const confirmPassword = document.getElementById("signup_confirmpass").value;
  let newUser = {
    name: username,
    email: email,
    password: password,
    confirmPassword: confirmPassword,
  };
  console.log(newUser);
  await fetch("/signup", {
    method: "POST",
    body: JSON.stringify(newUser),
  })
    .then((response) => response.json())
    .then((json) => {
      console.log(json);
      const modalHeading = document.getElementById("signup_result_heading");
      const modalBody = document.getElementById("signup_result_body");
      let modalBtn = document.getElementById("login");
      if (!json.status) {
        modalHeading.innerHTML = "Oh Snap!";
        modalBody.innerHTML = `${json.message}`;
        modalBtn.setAttribute("data-bs-target", "#signup_modal");
        modalBtn.textContent = "Sign up";
      } else {
        modalHeading.innerHTML = "Welcome!";
        modalBody.innerHTML = `You are now registered to Gritface!<br />
 You can now login.`;
        modalBtn.setAttribute("data-bs-target", "#login_modal");
        modalBtn.textContent = "Log in";
      }
    });
}
async function login() {
  const email = document.getElementById("login_email");
  const password = document.getElementById("login_pass");
  let user = {
    email: email.value,
    password: password.value,
  };
  password.value = "";
  resetLoginModal();
  console.log(user);
  await fetch("/login", {
    method: "POST",
    body: JSON.stringify(user),
  })
    .then((response) => response.json())
    .then((json) => {
      if (!json.status) {
        const loginPassLabel = document.getElementById("login_pass_label");
        loginPassLabel.innerHTML = `password<br />${json.message}`;
      } else {
        console.log(`logged in successfully with uid ${json.message}`);
        const loginForm = document.getElementById("login_success");
        loginForm.submit();
      }
    });
}
async function resetLoginModal() {
  const loginPassLabel = document.getElementById("login_pass_label");
  const loginPass = document.getElementById("login_pass");
  loginPassLabel.innerHTML = "password";
  loginPass.value = "";
  await fetch("/checkSession")
    .then((response) => response.json())
    .then((json) => {
      if (json.status) {
        window.location.replace("/loginSuccess");
      } else {
        console.log("please sign up");
      }
    });
}

async function loadPosts() {
  const lastPostID = parseInt(document.getElementById("container").lastChild.id)
  let lastPost = {
    lastPostID: lastPostID,
  };
  await fetch("/loadPosts", {
    method: "POST",
    body: JSON.stringify(lastPost),
  })
    .then((response) => response.json())
    .then((json) => {
      createPosts(json)
    });
}

function createPostDiv(postUserPic, postUsername, postID, postHeading, postBody, postCats, postInsertTime, postUpdateTime, commentsLength, comments, likeNum, dislikeNum) {
  return `<section class="row" id="post_section">
  <div class="row">
    <div class="col-2 col-md-1 col-lg-1 ms-2 mt-2">
      <img class="rounded-circle" style="max-width: 150%; border: 2px solid #54B4D3;" src="${postUserPic}">
    </div>
    <div class="col-7 mt-4">
      <h5 class="text-start text-info">${postUsername}</h5>
    </div>
  </div>
<div data-bs-target="#collapse_post_comments${postID}" data-bs-toggle="collapse">
    <div class="text-white rounded my-2 py-2" id="post_div">
        <div class="col-11 offset-1 my-1" id="post_heading">
            <h4>${postHeading}</h4>
        </div>
        <div class="col-10 offset-1" id="post_body">
            <div class="border-top border-info bg-dark text-center" id="post_image"></div>
            <div class="text-justify my-2">
                <pre>${postBody}</pre>
            </div>
            <div class="text-secondary">
            <p>${postCats}</p>
            <div class="row text-secondary">
                <div class="col-6 order-0 text-left" id="post_insert_time">
                    <p>${postInsertTime}</p>
                </div>
                <div class="col-6 order-1 text-end" id="post_mod_time">
                    <p>${postUpdateTime}</p>
                </div>
            </div>
        </div>
    </div>
</div>
</div>

<div class="offset-lg-1 offset-md-1 offset-0 py-1">
    <div class="col-12 mb-4 mb-lg-2 mb-md-2">
        <div class="row">
            <div class="mx-2" id="post_reactions_container${postID}">
                ${reactionButton(postID, 0, 1, likeNum)}
                ${reactionButton(postID, 0, 2, dislikeNum)}
                <p class="mx-1 text-info" id="number_of_comments"
                  data-bs-target="#collapse_post_comments${postID}" data-bs-toggle="collapse">
                ${commentsLength} Comment</p>
            </div>
        </div>
      </div>
      ${document.getElementById("user_name")? createCommentTextArea(postUserPic, postID):""}
  <div class="collapse" id="collapse_post_comments${postID}">
  ${comments}
  </div>
</section>`
}

function createCommentDiv(postID, commentID, commentUserPic, commentUsername, newComment, likeNumComment = 0, dislikeNumComment = 0) {
  return `
  <div class="row ms-2 pb-2" id="post_comments_container${postID}${commentID}">
    <div class="col-lg-9 offset-0 offset-md-1 offset-lg-1 col-md-10 col-11 border rounded bg-secondary" id="post_comment_body${postID}${commentID}">
      <div class="row pb-0 mb-0">
      <div class="col-md-1 col-lg-1 col-2 pt-1 me-4">
        <img class="rounded-circle" style="max-width: 50px; border: 2px solid #54B4D3" src="${commentUserPic}" id="user_pic">
        </div>
        <div class="col-5 pb-0 mb-0 h-50">
        <p class="text-info pt-1 mb-0 pb-0">${commentUsername}</p>
        </div>
        <pre class="pb-0 mb-0 offset-2"><p class="mb-0 pb-0" style="position:relative; top: -4px;">${newComment}</p></pre>
        </div>
  
  


  <div class="text-end pb-1 my-0" id="comment_reactions_container${postID}${commentID}">
  ${reactionButton(postID, commentID, 1, likeNumComment)}
  ${reactionButton(postID, commentID, 2, dislikeNumComment)}
  </div>
  </div>
  </div>`;
}

function createCommentTextArea(userPic, postID) {
  return `<div class="col-lg-10 col-md-10 col-12 justify-content-center ps-2 pe-2" id="user_comment">

  <!-- <div class="row">
  <div class="col-lg-2 col-md-2 d-none d-md-inline d-lg-inline">
  <img class="rounded-circle" style="max-width: 110%; border: 2px solid #54B4D3" src="${userPic}"></img>
  </div> -->

  <div class="col-12">
  <textarea
  id="newComment"
  class="bg-dark border-info rounded text-light pt-1 px-1"
  class="form-control"
  style="resize:none; font-size: 0.8em; width: 100%;"
  id="newComment"
  placeholder="Write a comment"></textarea>
  <button
  class="btn bg-info text-dark mb-1"
  type="button"
  style="font-size: 0.8em; float: right;"
  onclick="addComment(${postID})">
  Comment
  </button>
  </div>
  </div>
  </div>
  </div>`
};

function createPosts(json) {
  const container = document.getElementById("container");
  const userPic = document.getElementById("user_pic");
  for (const [key, postJSON] of Object.entries(json)) {
    // get categories
    let categories = "";
    if (postJSON.categories) {
      postJSON.categories.map(
        (category) => (categories += `#${category} `)
      );
    }
    // get like and dislike numbers
    let likeNum = 0,
      dislikeNum = 0;
    if (postJSON.reactions) {
      postJSON.reactions.map(function (reactions) {
        for (const [reaction_user_id, reaction] of Object.entries(
          reactions
        )) {
          if (reaction == "1") likeNum++;
          if (reaction == "2") dislikeNum++;
        }
      });
    }

    // create post div
    const postDiv = document.createElement("div");
    postDiv.classList.add("border", "rounded", "mx-auto", "col-lg-8", "col-md-8", "col-12", "mt-2", "mb-4", "mb-lg-2", "mb-md-2");
    postDiv.id = postJSON.post_id;

    // loop and create divs of comments
    let comments = ``;
    let likeNumComment, dislikeNumComment;
    for (const [key, comment] of Object.entries(postJSON.comments)) {
      if (comment.profile_image == "")
        comment.profile_image = "static/images/raccoon_thumbnail7.jpg";
      likeNumComment = 0;
      dislikeNumComment = 0;
      if (comment.reactions) {
        comment.reactions.map(function (reactions) {
          for (const [key, reaction] of Object.entries(reactions)) {
            if (reaction == "1") likeNumComment++;
            if (reaction == "2") dislikeNumComment++;
          }
        });
      }
      comments += createCommentDiv(postJSON.post_id, comment.comment_id, comment.profile_image, comment.username, comment.body, likeNumComment, dislikeNumComment);
    }
    
    // assemble the whole post div
    postDiv.innerHTML = createPostDiv(postJSON.profile_image, postJSON.username, postJSON.post_id, postJSON.heading, postJSON.body, categories, postJSON.insert_time, postJSON.update_time, Object.keys(postJSON.comments).length, comments, likeNum, dislikeNum);
    
    container.append(postDiv);
  }
}