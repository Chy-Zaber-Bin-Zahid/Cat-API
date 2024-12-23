document.addEventListener("DOMContentLoaded", function () {
  // Get DOM elements
  const upVoteBtn = document.querySelector(".vote-btn.up");
  const downVoteBtn = document.querySelector(".vote-btn.down");
  const catImagesDiv = document.getElementById("cat-images");
  const loading = document.querySelector(".loading-spinner");
  const navItems = document.querySelectorAll('.nav-item');
  const votingView = document.getElementById('voting-view');
  const breedsView = document.getElementById('breeds-view');

  // Initially hide breeds view
  breedsView.style.display = 'none';
  votingView.style.display = 'block';

  // Handle navigation
  navItems.forEach(item => {
    item.addEventListener('click', (e) => {
      e.preventDefault();
      const viewName = item.dataset.view;
      
      // Remove active class from all nav items
      navItems.forEach(nav => nav.classList.remove('active'));
      
      // Add active class to clicked nav items in both views
      document.querySelectorAll(`[data-view="${viewName}"]`)
        .forEach(nav => nav.classList.add('active'));
      
      // Show/hide views based on selection
      if (viewName === 'voting') {
        votingView.style.display = 'block';
        breedsView.style.display = 'none';
      } else if (viewName === 'breeds') {
        votingView.style.display = 'none';
        breedsView.style.display = 'block';
      }
    });
  });

  // Function to fetch new cat image
  async function fetchCatImage() {
    try {
      loading.style.display = "block";
      catImagesDiv.style.display = "none";
      
      const response = await fetch("/");
      const html = await response.text();
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, "text/html");
      const newImageURL = doc.querySelector("#cat-images img")?.src;

      if (newImageURL) {
        const newImageElement = document.createElement("img");
        newImageElement.src = newImageURL;
        newImageElement.alt = "Cat image";
        newImageElement.classList.add("pet-image");

        catImagesDiv.innerHTML = "";
        catImagesDiv.appendChild(newImageElement);
        catImagesDiv.style.display = "block";
        loading.style.display = "none";
      } else {
        throw new Error("No image URL found");
      }
    } catch (error) {
      console.error("Error fetching new cat image:", error);
      loading.innerHTML = "<p>Error loading image. Please try again.</p>";
    }
  }

  // Handle voting buttons
  function handleVote() {
    fetchCatImage();
  }

  upVoteBtn?.addEventListener("click", handleVote);
  downVoteBtn?.addEventListener("click", handleVote);
});