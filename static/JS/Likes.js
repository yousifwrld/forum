function handleReaction(id, action, type) {
  fetch("/reaction", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ ID: id, action: action, type: type }),
  })
    .then((response) => response.json())
    .then((data) => {
      console.log("Received data:", data); // Log the received data

      if (data.likes !== undefined && data.dislikes !== undefined) {
        document.getElementById("likeCount-" + id).textContent = data.likes;
        document.getElementById("dislikeCount-" + id).textContent =
          data.dislikes;

        // Update the like icon
        const likeIcon = document.getElementById("likeIcon-" + id);
        console.log("user_liked:", data.user_liked); // Log user_liked status
        if (data.user_liked) {
          likeIcon.classList.remove("bx-like");
          likeIcon.classList.add("bxs-like");
        } else {
          likeIcon.classList.remove("bxs-like");
          likeIcon.classList.add("bx-like");
        }

        // Update the dislike icon
        const dislikeIcon = document.getElementById("dislikeIcon-" + id);
        console.log("user_disliked:", data.user_disliked); // Log user_disliked status
        if (data.user_disliked) {
          dislikeIcon.classList.remove("bx-dislike");
          dislikeIcon.classList.add("bxs-dislike");
        } else {
          dislikeIcon.classList.remove("bxs-dislike");
          dislikeIcon.classList.add("bx-dislike");
        }
      } else {
        console.error("Invalid response data:", data);
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      window.location.replace("/login");
    });
}
