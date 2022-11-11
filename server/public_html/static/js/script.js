async function initPage() {
    await fetch("/posts")
      .then((response) => response.json())
      .then(function (json) {
        console.log(json)
        for (const [key, postJSON] of Object.entries(json)) {

          const postSection = document.getElementById("post_section");
  
          const postDiv = document.createElement("div");
          const postBody = document.createElement("div");
          const postHeading = document.createElement("div");
          const postInsertTime = document.createElement("div");
          const postModTime = document.createElement("div");
          const postImage = document.createElement("div");
          const postComments = document.createElement("div");
          const postReactions = document.createElement("div");
          const postRow = document.createElement("div");
          const postLike = document.createElement("button");
          const postHeart = document.createElement("button");
          const postComment = document.createElement("div");
          const postCommentImg = document.createElement("div");
  
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
          postHeading.classList.add("col-10", "offset-1", "my-2");
          postInsertTime.classList.add("col-6", "order-0", "text-left");
          postModTime.classList.add("col-6", "order-1", "text-end");
          postComments.classList.add("row", "my-2");
          postReactions.classList.add("col-5", "mx-1");
          postImage.classList.add("border", "bg-info", "text-center");
          postLike.classList.add("btn", "border", "rounded", "bg-dark");
          postHeart.classList.add("btn", "border", "rounded", "bg-dark");
          postComment.classList.add("col-9", "border", "rounded", "bg-secondary");
          postCommentImg.classList.add("col-2", "mx-1", "border", "rounded", "bg-info");
  
          postBody.textContent = postJSON.body;
          postHeading.textContent = postJSON.heading;
          postInsertTime.textContent = postJSON.insert_time;
          postModTime.textContent = postJSON.update_time;
          postComments.textContent = JSON.stringify(postJSON.comments);
          postReactions.textContent = postJSON.post_reactions;
          postImage.textContent = postJSON.image;
  
          postSection.appendChild(postDiv);
          postDiv.appendChild(postHeading);
          postDiv.appendChild(postBody);
          postDiv.appendChild(postRow)
          postBody.appendChild(postImage);
          postBody.appendChild(postInsertTime);
          postBody.appendChild(postModTime);
          postRow.appendChild(postReactions);
          postRow.appendChild(postComments);
          postReactions.appendChild(postLike);
          postReactions.appendChild(postHeart);
          console.log(postJSON)
        };
      });
  }
  