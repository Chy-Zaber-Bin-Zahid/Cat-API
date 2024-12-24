document.addEventListener("DOMContentLoaded", function () {
  // Get DOM elements
  const upVoteBtn = document.querySelector(".vote-btn.up");
  const downVoteBtn = document.querySelector(".vote-btn.down");
  const catImagesDiv = document.getElementById("cat-images");
  const loading = document.querySelector(".loading-spinner");
  const navItems = document.querySelectorAll('.nav-item');
  const votingView = document.getElementById('voting-view');
  const breedsView = document.getElementById('breeds-view');
  const favsView = document.getElementById('favs-view');
  const viewToggles = document.querySelectorAll('.view-toggle');

//   document.addEventListener('DOMContentLoaded', function() {
//     // Set the default breed ID (you can set this to whatever you want)
//     var defaultBreedId = 'acur';  // Example breed ID, you can change it to a valid ID
//     console.log("Selected Breed ID:", defaultBreedId);
//     // Set the default value in the dropdown
//     var breedSelect = document.getElementById('breedSelect');
//     breedSelect.value = defaultBreedId;

//     // Trigger the change event manually to load images for the default breed
//     breedSelect.dispatchEvent(new Event('change'));
// });

// Define the function with a name
async function loadBreedImages() {
  const selectedId = document.getElementById('breedSelect').value; // Get the selected option's value (ID)
  console.log("Id: ", selectedId);

  try {
      const response = await fetch(`/catImages?breed_id=${selectedId}`, {
          method: 'GET',
      });

      if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const images = await response.json();

      // Get the carousel container
      const carouselContainer = document.querySelector('.carousel-slide');
      carouselContainer.innerHTML = ''; // Clear existing content

      // Loop through the images and add them to the carousel
      images.forEach(image => {
          const imgElement = document.createElement('img');
          const dotElement = document.createElement('div');
          imgElement.src = image.url;
          imgElement.alt = "Breed Image";
          imgElement.className = "carousel-image"; // Optional: Add a class for styling
          dotElement.className = "dot"

          carouselContainer.appendChild(imgElement);
          document.querySelector('.carousel-dots').appendChild(dotElement)
      });

      console.log("Images loaded successfully!");
  } catch (error) {
      console.error('Error fetching data:', error);
  }
}

// Add the event listener
document.getElementById('breedSelect').addEventListener('change', loadBreedImages);





  // Initially hide breeds view
  breedsView.style.display = 'none';
  favsView.style.display = 'none';
  votingView.style.display = 'block';

  viewToggles.forEach(item => {
    item.addEventListener('click', (e) => {
      e.preventDefault();
      const viewClicked = item.dataset.view;
      viewToggles.forEach(nav => nav.classList.remove('active'));
      document.querySelectorAll(`[data-view="${viewClicked}"]`)
        .forEach(nav => nav.classList.add('active'));
      });
    });
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
        favsView.style.display = 'none';
      } else if (viewName === 'breeds') {
        votingView.style.display = 'none';
        breedsView.style.display = 'block';
        favsView.style.display = 'none';
        loadBreedImages()
      } else {
        votingView.style.display = 'none';
        breedsView.style.display = 'none';
        favsView.style.display = 'block';
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