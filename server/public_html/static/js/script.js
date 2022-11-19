window.onload = initPage();

async function initPage() {
  await fetch("/posts")
    .then((response) => response.json())
    .then(function (json) {
      for (const [key, postJSON] of Object.entries(json)) {
        const postDiv = document.createElement("div");
        postDiv.classList.add(
          "border",
          "rounded",
          "content",
          "mx-auto",
          "col-8"
        );
        let comments = "";

        for (const [key, comment] of Object.entries(postJSON.comments)) {
          comments += `<div class="collapse" id="collapse_post_comments">
            <div class="row my-3" id="post_comments">
                <div class="col-1 row-1 mx-1 border rounded-start bg-info">${comment.user_id}</div>
                <div class="col-8 border rounded-end bg-secondary" id="post_comments">
                ${comment.body}
                </div>
            </div>
        </div>`;
          postDiv.innerHTML = `<section class="row" id="post_section">
        <div data-bs-target="#collapse_post_comments" data-bs-toggle="collapse">
            <div class="text-white rounded my-2 py-2" id="post_div">
                <div class="col-11 offset-1 my-1" id="post_heading">
                    ${postJSON.heading}
                </div>
                <div class="col-10 offset-1" id="post_body">
                    <div class="border bg-info text-center" id="post_image">${postJSON.image}"</div>
                    <div class="text-justify my-2">
                        ${postJSON.body}
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
                    <div class="mx-1" id="post_reactions">
                        <button class="bg-dark border rounded-start">👍<span
                                class="badge text-info">10</span></button>
                        <button class="bg-dark border">👎<span class="badge text-info">5</span></button>
                        <button class="bg-dark border rounded-end">💛<span
                                class="badge text-info">8</span></button>
                    </div>
                </div>
            </div>
            ${comments}
        </div>
    </section>`;
        }

        // loop and create comments
        const container = document.getElementById("container");
        container.append(postDiv);
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
