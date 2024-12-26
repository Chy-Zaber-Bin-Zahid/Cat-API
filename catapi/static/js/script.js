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
  const likeButton = document.querySelector('.like-btn');
  const viewToggles = document.querySelectorAll('.view-toggle');

// Define the function with a name
async function getFavourites() {
  const apiKey = 'live_Ii20w7Wt785t9kCsxDQYAMTIIL7epsK1IaGiHL3hxWw0ou2AfkvZ3FAMxJ4NEc0Z'; // Replace with your actual API key

  try {
      // Send GET request to TheCatAPI to fetch favourites
      const response = await fetch('https://api.thecatapi.com/v1/favourites', {
          method: 'GET',
          headers: {
              'x-api-key': apiKey // Send API key in the header
          }
      });

      if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();

      // Get the container where we want to display the favourite images
      const favsGrid = document.querySelector('.favs-grid');
      favsGrid.innerHTML = ''; // Clear existing content

      // Loop through the favourites and create image elements
      data.forEach(fav => {
          const imgElement = document.createElement('img');
          imgElement.src = fav.image.url;
          imgElement.alt = 'Favourite cat image';
          imgElement.className = 'fav-image'; // You can style the images with this class

          // Append the image to the favs-grid container
          favsGrid.appendChild(imgElement);
      });

      console.log("Favourites loaded successfully!");

  } catch (error) {
      console.error('Error fetching favourites:', error);
  }
}


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

      const data = await response.json();

      // Extract images and breed info
      const { images, breedInfo } = data;

      // Get the carousel container
      const carouselContainer = document.querySelector('.carousel-slide');
      carouselContainer.innerHTML = ''; // Clear existing content

      // Reset the slide view to the first image
      carouselContainer.style.transform = 'translateX(0px)';

      // Loop through the images and add them to the carousel
      let count = 0;

      images.forEach((image, index) => {
          const imgElement = document.createElement('img');
          const dotElement = document.createElement('div');

          // Set image attributes
          imgElement.src = image.url;
          imgElement.alt = `Slide ${count + 1}`;
          imgElement.className = "carousel-image"; // Optional: for styling purposes

          // Set dot class
          dotElement.className = "dot";

          // Append the image to the carousel container
          carouselContainer.appendChild(imgElement);

          // Select the carousel-dots element
          const carouselDots = document.querySelector('.carousel-dots');

          // Remove all existing child elements from carousel-dots only once
          if (count === 0) {
              while (carouselDots.firstChild) {
                  carouselDots.removeChild(carouselDots.firstChild);
              }
          }

          // Add the 'active' class to the first dot
          if (count === 0) dotElement.classList.add('active'); // Make the first dot active
          dotElement.dataset.index = index;

          // Append the dot to the carousel-dots element
          carouselDots.appendChild(dotElement);

          count++;
      });

      // Handle dot clicks
      const dots = document.querySelectorAll('.dot');

      dots.forEach(dot => {
          dot.addEventListener('click', () => {
              // Get the index of the clicked dot
              const index = parseInt(dot.dataset.index);

              // Update active dot
              dots.forEach(d => d.classList.remove('active'));
              dot.classList.add('active');

              // Scroll the corresponding image into view
              const imageWidth = carouselContainer.querySelector('.carousel-image').clientWidth;
              carouselContainer.style.transform = `translateX(-${index * imageWidth}px)`;
          });
      });

      // Update breed details
      const breedName = document.getElementById('breedName');
      const breedOrigin = document.getElementById('breedOrigin');
      const breedDesc = document.getElementById('breedDesc');
      const wikiLink = document.getElementById('wikiLink');

      // Clear and set the name and origin properly
      breedName.textContent = breedInfo.name;
      breedOrigin.textContent = ` (${breedInfo.origin})`;
      breedName.appendChild(breedOrigin); // Ensure the span is inside the heading

      // Update the description
      breedDesc.textContent = breedInfo.description;

      // Update the Wikipedia link
      wikiLink.href = breedInfo.wikipedia_url;
      wikiLink.textContent = "Learn More on Wikipedia";

      console.log("Images and breed info loaded successfully!");
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

  // viewToggles.forEach(item => {
  //   item.addEventListener('click', (e) => {
  //     e.preventDefault();
  //     const viewClicked = item.dataset.view;
  //     viewToggles.forEach(nav => nav.classList.remove('active'));
  //     document.querySelectorAll(`[data-view="${viewClicked}"]`)
  //       .forEach(nav => nav.classList.add('active'));
  //     });
  //   });
    const favsGrid = document.querySelector(".favs-grid");

    viewToggles.forEach((toggle) => {
      toggle.addEventListener("click", () => {
        const view = toggle.getAttribute("data-view");

        // Remove 'active' class from all buttons
        viewToggles.forEach((btn) => btn.classList.remove("active"));

        // Add 'active' class to the clicked button
        toggle.classList.add("active");

        // Toggle the 'list-view' class on the favs-grid
        if (view === "list") {
          favsGrid.classList.add("list-view");
        } else {
          favsGrid.classList.remove("list-view");
        }
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
        getFavourites()
      }
    });
  });

  // Function to fetch new cat image
let currentImageID = null; // To store the current image ID globally

async function fetchCatImage(name) {
  try {
    if (name !== "like") {
      console.log("Fetching new cat image...");
      // Fetch a new image
      loading.style.display = "block";
      catImagesDiv.style.display = "none";

      const response = await fetch("/");
      const html = await response.text();
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, "text/html");

      const newImageURL = doc.querySelector("#cat-images img")?.src;
      const newImageID = doc.querySelector("#cat-images img")?.id;

      if (newImageURL && newImageID) {
        // Store the current image details globally
        currentImageID = newImageID;
        currentImageURL = newImageURL;

        const newImageElement = document.createElement("img");
        newImageElement.src = newImageURL;
        newImageElement.alt = "Cat image";
        newImageElement.classList.add("pet-image");

        catImagesDiv.innerHTML = "";
        catImagesDiv.appendChild(newImageElement);
        catImagesDiv.style.display = "block";
        loading.style.display = "none";
      } else {
        throw new Error("No image URL or ID found");
      }
    } else {
      // Add the current image to favorites
      if (currentImageID) {
        console.log("Cat Image ID:", currentImageID);

        const rawBody = JSON.stringify({
          image_id: currentImageID, 
          sub_id: "user-123",
        });

        const newFavourite = await fetch("https://api.thecatapi.com/v1/favourites", {
          method: "POST",
          headers: {
            "x-api-key": "live_Ii20w7Wt785t9kCsxDQYAMTIIL7epsK1IaGiHL3hxWw0ou2AfkvZ3FAMxJ4NEc0Z", // Replace with your actual API key
            "Content-Type": "application/json",
          },
          body: rawBody,
        });

        const result = await newFavourite.json();
        console.log(result); // Handle the response as needed
      } else {
        console.error("No current image ID found to add to favorites.");
      }
      // Fetch a new image
      loading.style.display = "block";
      catImagesDiv.style.display = "none";

      const response = await fetch("/");
      const html = await response.text();
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, "text/html");

      const newImageURL = doc.querySelector("#cat-images img")?.src;
      const newImageID = doc.querySelector("#cat-images img")?.id;

      if (newImageURL && newImageID) {
        // Store the current image details globally
        currentImageID = newImageID;
        currentImageURL = newImageURL;

        const newImageElement = document.createElement("img");
        newImageElement.src = newImageURL;
        newImageElement.alt = "Cat image";
        newImageElement.classList.add("pet-image");

        catImagesDiv.innerHTML = "";
        catImagesDiv.appendChild(newImageElement);
        catImagesDiv.style.display = "block";
        loading.style.display = "none";
      } else {
        throw new Error("No image URL or ID found");
      }
    }
  } catch (error) {
    console.error("Error fetching new cat image:", error);
    loading.innerHTML = "<p>Error loading image. Please try again.</p>";
  }
}

window.onload = () => {
  fetchCatImage("reload"); // Call the function without "like" to fetch a new image
};

  // Handle voting buttons
  function handleVote(name) {
    fetchCatImage(name);
  }

  upVoteBtn?.addEventListener("click", () => handleVote("up"));
  downVoteBtn?.addEventListener("click", () => handleVote("down"));
  likeButton?.addEventListener("click", () => handleVote("like"));
  
});
