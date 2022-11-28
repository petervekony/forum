function userPost() {
    console.log("userPost");
    window.onload = filterPage()
  }
  
  

async function filterPage() {
    await fetch("/userPosts")
    .then((response) => response.json())
    .then(function (json) {
      console.log(json);
      let commentTextArea = "";
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
              if (reaction == "⬆️") likeNum++;
              if (reaction == "⬇️") dislikeNum++;
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
)}">
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
          likeNumComment = 0;
          dislikeNumComment = 0;
          if (comment.reactions) {
            console.log(comment.reactions);
            comment.reactions.map(function (reactions) {
              for (const [key, reaction] of Object.entries(reactions)) {
                if (reaction == "⬆️") likeNumComment++;
                if (reaction == "⬇️") dislikeNumComment++;
              }
            });
          }
          comments += `
 <div class="row my-3 ms-auto" id="post_comments">
 <div class="col-1 mx-2">
 <img class="rounded-circle" style="max-width: 120%; border: 2px solid #54B4D3" src="${userPic.getAttribute(
  "src"
)}">
 </div>
 <div class="col-8 border rounded bg-secondary" id="post_comments">
 <p class="text-info pt-1">${comment.username}</p>
 <pre><p>${comment.body}</p></pre>
 <div class="row">
 <div class="text-end mb-1" id="comment_reactions">
 <button class="btn btn-dark border">⬆️
 <span class="badge text-info">${likeNumComment}</span>
 </button>
 <button class="btn btn-dark border">⬇️
 <span class="badge text-info">${dislikeNumComment}</span>
 </button>
 </div>
 </div>
 </div>
 </div>`;
        }
        comments += "</div>";
        postDiv.innerHTML = `<section class="row" id="post_section">
        <h5 class="text-start mx-3 mt-2 text-info">${postJSON.username}</h5>
        <div data-bs-target="#collapse_post_comments${
          postJSON.post_id
        }" data-bs-toggle="collapse">
            <div class="text-white rounded my-2 py-2" id="post_div">
                <div class="col-11 offset-1 my-1" id="post_heading">
                    <h4>${postJSON.heading}</h4>
                </div>
                <div class="col-10 offset-1" id="post_body">
                    <div class="border bg-info text-center" id="post_image">${
                      postJSON.image
                    }"</div>
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
                  <div class="mx-1" id="post_reactions">
                 <button class="btn btn-dark border">⬆️<span
                  class="badge text-info">${likeNum}</span></button>
                  <button class="btn btn-dark border">⬇️<span class="badge text-info">${dislikeNum}</span></button>
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
        const container = document.getElementById("container");
        container.append(postDiv);
      }
    });
}
