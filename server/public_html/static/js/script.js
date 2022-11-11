async function initPage() {
    await fetch("/posts")
      .then((response) => response.json())
      .then(function (json) {
        json.map(function (postJSON) {
          const postSection = document.getElementById("post_section");
  
          const postDiv = document.createElement("div");
          const postBody = document.createElement("div");
          const postHeading = document.createElement("div");
          const postInsertTime = document.createElement("div");
          const postModTime = document.createElement("div");
          const postImage = document.createElement("div");
          const postComments = document.createElement("div");
          const postReactions = document.createElement("div");
  
          //test
          const a = document.createElement("div");
          const b = document.createElement("div");
          const c = document.createElement("div");
          const d = document.createElement("div");
  
          postDiv.classList.add(
            "col-8",
            "offset-2",
            "text-center",
            "border",
            "rounded",
            "my-2",
            "py-2"
          );
          postBody.classList.add("postBody")
          postHeading.classList.add("postHeading")
          postInsertTime.classList.add("postInsertTime")
          postModTime.classList.add("postModTime")
          postComments.classList.add("postComments")
          postReactions.classList.add("postReactions")
          postImage.classList.add("postImage")
  
          postBody.textContent = postJSON.post_body;
          postHeading.textContent = postJSON.post_heading;
          postInsertTime.textContent = postJSON.post_insert_time;
          postModTime.textContent = postJSON.post_mod_time;
          postComments.textContent = postJSON.post_comments;
          postReactions.textContent = postJSON.post_reactions;
          postImage.textContent = postJSON.post_image;
  
          postSection.appendChild(postDiv);
          postDiv.appendChild(postHeading);
          postDiv.appendChild(postBody);
          postDiv.appendChild(postReactions);
          postDiv.appendChild(postComments);
          postBody.appendChild(postImage);
          postBody.appendChild(postInsertTime);
          postBody.appendChild(postModTime);
        });
      });
  }
  