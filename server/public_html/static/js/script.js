window.onload = initPage();

async function initPage(request = "/posts") {
  console.log("hello from initpage, " + request);
  await fetch(request)
    .then((response) => response.json())
    .then(function (json) {
      let commentTextArea = "";
      const container = document.getElementById("container");
      container.innerHTML = "";
      const userPic = document.getElementById("user_pic");
      for (const [key, postJSON] of Object.entries(json)) {
        let categories = "";
        if (postJSON.categories) {
          postJSON.categories.map(
            (category) => (categories += `#${category} `)
          );
        }
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
        const postDiv = document.createElement("div");
        postDiv.classList.add("border", "rounded", "mx-auto", "col-8", "mt-2");
        postDiv.id = postJSON.post_id;
        let comments = `<div class="collapse" id="collapse_post_comments${postJSON.post_id}">`;
        if (document.getElementById("user_name")) {
          commentTextArea = `<div class="col-10 justify-content-center mx-2 mb-2" id="user_comment">
 <div class="row">
 <div class="col-1 mx-2">
 <img class="rounded-circle" style="max-width: 150%; border: 2px solid #54B4D3" src="${userPic.getAttribute(
   "src"
 )}"></img>
 </div>
 <div class="col-10 text-start">
 <div class="input-group">
 <textarea
 id="newComment"
 class="bg-dark border-info rounded text-light px-2 w-75"
 class="form-control"
 style="resize:none;"
 id="newComment"
 placeholder="Write a comment"></textarea>
 <div class="input-group-append mx-2">
 <button
 class="btn bg-info text-dark mt-2"
 type="button"
 onclick="addComment(${postJSON.post_id})">
 Comment
 </button>
 </div>
 </div>
 </div>
 </div>
 </div>`;
        }
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
          // if (comment.profile_image == "") comment.profile_image = "static/images/raccoon_thumbnail7.jpg";
          comments += `
 <div class="row ms-auto pb-1" id="post_comments_container${postJSON.post_id}${comment.comment_id}">
 <div class="col-1 mx-2">
 <img class="rounded-circle" style="max-width: 120%; border: 2px solid #54B4D3" src="${comment.profile_image}" id="user_pic">
 </div>
 <div class="col-8 border rounded bg-secondary" id="post_comment_body${postJSON.post_id}${comment.comment_id}">
 <p class="text-info pt-1">${comment.username}</p>
 <pre><p>${comment.body}</p></pre>


 <div class="text-end pb-1 my-0" id="comment_reactions_container${postJSON.post_id}${comment.comment_id}">
 ${reactionButton(postJSON.post_id, comment.comment_id, 1, likeNumComment)}
 ${reactionButton(postJSON.post_id, comment.comment_id, 2, dislikeNumComment)}
 </div>
 </div>
 </div>`;
        }
        comments += "</div>";
        postDiv.innerHTML = `<section class="row" id="post_section">
        <div class="row">
          <div class="col-1 ms-2 mt-2">
          <img class="rounded-circle" style="max-width: 120%; border: 2px solid #54B4D3" src="${
            postJSON.profile_image
          }" id="user_pic">
          </div>
          <div class="col-7 mt-4">
          <h5 class="text-start text-info">${postJSON.username}</h5>
          </div>
        </div>
        
        <div data-bs-target="#collapse_post_comments${
          postJSON.post_id
        }" data-bs-toggle="collapse">
            <div class="text-white rounded my-2 py-2" id="post_div">
                <div class="col-11 offset-1 my-1" id="post_heading">
                    <h4>${postJSON.heading}</h4>
                </div>
                <div class="col-10 offset-1" id="post_body">
                    <div class="border-top bg-dark border-info text-center" id="post_image">
                    </div>
                    <div class="text-justify my-2">
                        <pre><p>${postJSON.body}</p></pre>
                    </div>
                    <div class="text-secondary">
                    <p>${categories}</p>
                    </div>
                    <div class="row text-secondary">
                        <div class="col-6 order-0 text-left" id="post_insert_time">
                            ${postJSON.insert_time}
                        </div>
                        <div class="col-6 order-1 text-end" id="post_mod_time">
                            ${postJSON.update_time}
                        </div>
                      </div>
                    </div>
                  </div>
                  </div>
                   <div class="offset-1 py-1">
                   <div class="col-12 mb-2">
                  <div class="row">

                  <div class="mx-1">

                  ${reactionButton(postJSON.post_id, 0, 1, likeNum)}
                  
                  ${reactionButton(postJSON.post_id, 0, 2, dislikeNum)}


                    <p class="mx-1 text-info" id="number_of_comments">${
                      Object.keys(postJSON.comments).length
                    } Comments</p>
                  </div>
                  </div>

  
             ${commentTextArea}
            ${comments}
               </div>
          </section>`;

        // loop and create comments

        container.append(postDiv);
      }
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

  returndata += `<button class="btn btn-dark ${addClass}" style="${setStyle}" id="rbc${postId}${commentId}${reactId}" onclick="addReaction(${postId}, ${commentId}, ${reactId}, this)">${reactIcon}
                  <span class="badge text-info" id="rb${postId}${commentId}${reactId}">${reactCount}</span>
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
      for(let i=1; i < 3; i++) {
        document.getElementById("rb" + postID + commentID + i).innerHTML = json['rb'+i];
        document.getElementById("rbc" + postID + commentID + i).classList.remove("active")
      }
      targetButton.classList.add("active");
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
        // const uid = document.getElementById("login_email");
        // uid.value = json.message;
        // test
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
      let commentTextArea = "";
      const container = document.getElementById("container");
      const userPic = document.getElementById("user_pic");
      for (const [key, postJSON] of Object.entries(json)) {
        let categories = "";
        if (postJSON.categories) {
          postJSON.categories.map(
            (category) => (categories += `#${category} `)
          );
        }
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
        const postDiv = document.createElement("div");
        postDiv.classList.add("border", "rounded", "mx-auto", "col-8", "mt-2");
        postDiv.id = postJSON.post_id;
        let comments = `<div class="collapse" id="collapse_post_comments${postJSON.post_id}">`;
        if (document.getElementById("user_name")) {
          commentTextArea = `<div class="col-10 justify-content-center mx-2 mb-2" id="user_comment">
 <div class="row">
 <div class="col-1 mx-2">
 <img class="rounded-circle" style="max-width: 150%; border: 2px solid #54B4D3" src="${userPic.getAttribute(
   "src"
 )}"></img>
 </div>
 <div class="col-10 text-start">
 <div class="input-group">
 <textarea
 id="newComment"
 class="bg-dark border-info rounded text-light px-2 w-75"
 class="form-control"
 style="resize:none;"
 id="newComment"
 placeholder="Write a comment"></textarea>
 <div class="input-group-append mx-2">
 <button
 class="btn bg-info text-dark mt-2"
 type="button"
 onclick="addComment(${postJSON.post_id})">
 Comment
 </button>
 </div>
 </div>
 </div>
 </div>
 </div>`;
        }
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
          // if (comment.profile_image == "") comment.profile_image = "static/images/raccoon_thumbnail7.jpg";
          comments += `
 <div class="row ms-auto" id="post_comments">
 <div class="col-1 mx-2">
 <img class="rounded-circle" style="max-width: 120%; border: 2px solid #54B4D3" src="${comment.profile_image}" id="user_pic">
 </div>
 <div class="col-8 border rounded bg-secondary" id="post_comments">
 <p class="text-info pt-1">${comment.username}</p>
 <pre><p>${comment.body}</p></pre>
 <div class="row">


 <div class="text-end mb-1" id="comment_reactions">
 ${reactionButton(postJSON.post_id, comment.comment_id, 1, likeNumComment)}
 ${reactionButton(postJSON.post_id, comment.comment_id, 2, dislikeNumComment)}
 </div>
 </div>
 </div>`;
        }
        comments += "</div>";
        postDiv.innerHTML = `<section class="row" id="post_section">
        <div class="row">
          <div class="col-1 ms-2 mt-2">
          <img class="rounded-circle" style="max-width: 120%; border: 2px solid #54B4D3" src="${
            postJSON.profile_image
          }" id="user_pic">
          </div>
          <div class="col-7 mt-4">
          <h5 class="text-start text-info">${postJSON.username}</h5>
          </div>
        </div>
        
        <div data-bs-target="#collapse_post_comments${
          postJSON.post_id
        }" data-bs-toggle="collapse">
            <div class="text-white rounded my-2 py-2" id="post_div">
                <div class="col-11 offset-1 my-1" id="post_heading">
                    <h4>${postJSON.heading}</h4>
                </div>
                <div class="col-10 offset-1" id="post_body">
                    <div class="border-top bg-dark border-info text-center" id="post_image">
                    </div>
                    <div class="text-justify my-2">
                        <pre><p>${postJSON.body}</p></pre>
                    </div>
                    <div class="text-secondary">
                    <p>${categories}</p>
                    </div>
                    <div class="row text-secondary">
                        <div class="col-6 order-0 text-left" id="post_insert_time">
                            ${postJSON.insert_time}
                        </div>
                        <div class="col-6 order-1 text-end" id="post_mod_time">
                            ${postJSON.update_time}
                        </div>
                      </div>
                    </div>
                  </div>
                  </div>
                   <div class="offset-1 py-1">
                   <div class="col-12 mb-2">
                  <div class="row">

                  <div class="mx-1">
                  
                  ${reactionButton(postJSON.post_id, 0, 1, likeNum)}
                  
                  ${reactionButton(postJSON.post_id, 0, 2, dislikeNum)}



                    <p class="mx-1 text-info" id="number_of_comments">${
                      Object.keys(postJSON.comments).length
                    } Comments</p>
                  </div>
                  </div>

  
             ${commentTextArea}
            ${comments}
               </div>
          </section>`;

        // loop and create comments

        container.append(postDiv);
      }
    });
}