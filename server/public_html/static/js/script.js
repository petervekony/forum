window.onload = initPage();

async function initPage() {
  await fetch("/posts")
    .then((response) => response.json())
    .then(function (json) {
      for (const [key, postJSON] of Object.entries(json)) {
        const postSection = document.getElementById("post_section");

        const postDiv = document.createElement("div");
        postDiv.id = postJSON.post_id;

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

        postBodyText.textContent = postJSON.body;
        postHeading.textContent = postJSON.heading;
        postInsertTime.textContent = postJSON.insert_time;
        postModTime.textContent = postJSON.update_time;
        postReactions.textContent = postJSON.post_reactions;
        postImage.textContent = `<img src="/server/public_html/statis/images/${postJSON.image}">`;
        postLike.textContent = "👍";
        postDislike.textContent = "👎";
        postHeart.textContent = "❤️";
        // postLikeNum.textContent = "0";
        // test

        postSection.appendChild(postDiv);
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

        // loop and create comments
        for (const [key, comment] of Object.entries(postJSON.comments)) {
          const postCommentRow = document.createElement("div");
          const postComment = document.createElement("div");
          const postCommenterImg = document.createElement("div");
          postCommentRow.classList.add("row", "my-2");
          postComment.classList.add(
            "col-9",
            "border",
            "rounded",
            "bg-secondary"
          );
          postCommenterImg.classList.add(
            "col-2",
            "mx-1",
            "border",
            "rounded",
            "bg-info"
          );
          postComment.textContent = comment.body;
          postCommenterImg.textContent = `User id:${comment.user_id} profile pic`;
          postCommentRow.appendChild(postCommenterImg);
          postCommentRow.appendChild(postComment);
          postBody.appendChild(postCommentRow);
        }
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
      const modalBtn = document.getElementById("login");

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
