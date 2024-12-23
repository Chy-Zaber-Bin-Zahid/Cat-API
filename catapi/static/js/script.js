document.addEventListener("DOMContentLoaded", function () {
    const upVoteBtn = document.querySelector(".vote-btn.up");
    const downVoteBtn = document.querySelector(".vote-btn.down");
    const catImagesDiv = document.getElementById("cat-images");
    const loading = document.querySelector(".loading-spinner");
  
    // Function to fetch a new cat image from the server
    function fetchCatImage() {
      fetch("/")  // Send a GET request to the root URL to get a new image
        .then(response => response.text())  // Get the HTML response
        .then(html => {
          const parser = new DOMParser();
          const doc = parser.parseFromString(html, "text/html");
          const newImageURL = doc.querySelector("#cat-images img")?.src;
  
          if (newImageURL) {
            // Update the image source with the new URL
            const newImageElement = document.createElement("img");
            newImageElement.src = newImageURL;
            newImageElement.alt = "New Cat Image";
            newImageElement.classList.add("pet-image");
  
            // Replace the old image with the new one
            catImagesDiv.innerHTML = "";  // Clear the current image
            catImagesDiv.appendChild(newImageElement);  // Add the new image
            catImagesDiv.style.display = "block"
            loading.style.display = "none"
          } else {
            console.log("No image URL found.");
          }
        })
        .catch(error => {
          console.error("Error fetching new cat image:", error);
        });
    }
  
    // Attach event listeners to the vote buttons
    upVoteBtn.addEventListener("click", function () {
      console.log("Up vote clicked");
      catImagesDiv.style.display = "none"
      loading.style.display = "block"
      fetchCatImage();
    });
  
    downVoteBtn.addEventListener("click", function () {
      console.log("Down vote clicked");
      catImagesDiv.style.display = "none"
      loading.style.display = "block"
      fetchCatImage();
    });
  });
  