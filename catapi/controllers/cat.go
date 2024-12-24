package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

// Struct to parse the breed data from the API
type Breed struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Desc   string `json:"description"`
	Wiki   string `json:"wikipedia_url"`
}

// Struct to parse the cat image data
type CatImage struct {
	URL string `json:"url"`
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	// Channel to get the breed and image data
	breedChan := make(chan []Breed)
	imageChan := make(chan string)

	// Start a goroutine to fetch the breed data
	go func() {
		apiURL := "https://api.thecatapi.com/v1/breeds"
		resp, err := http.Get(apiURL)
		if err != nil {
			// If there's an error with the HTTP request
			fmt.Println("Error fetching breed data:", err)
			breedChan <- nil
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// If there's an error reading the response body
			fmt.Println("Error reading response body:", err)
			breedChan <- nil
			return
		}

		var breeds []Breed
		err = json.Unmarshal(body, &breeds)
		if err != nil {
			// If the JSON parsing fails
			fmt.Println("Error parsing JSON:", err)
			breedChan <- nil
			return
		}

		// Send the breed data to the channel
		breedChan <- breeds
	}()

	// Start a goroutine to fetch the cat image
	go func() {
		apiURL := "https://api.thecatapi.com/v1/images/search"
		resp, err := http.Get(apiURL)
		if err != nil {
			// If there's an error with the HTTP request
			fmt.Println("Error fetching cat image:", err)
			imageChan <- ""
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// If there's an error reading the response body
			fmt.Println("Error reading response body:", err)
			imageChan <- ""
			return
		}

		var catImages []CatImage
		err = json.Unmarshal(body, &catImages)
		if err != nil || len(catImages) == 0 {
			// If the JSON parsing fails or the array is empty
			fmt.Println("Error parsing JSON:", err)
			imageChan <- ""
			return
		}

		// Send the URL of the first image to the channel
		imageChan <- catImages[0].URL
	}()

	// Fetch breed and image data from the channels
	breeds := <-breedChan
	catImageURL := <-imageChan

	if breeds == nil {
		// If there's an error (i.e., the breed data is nil)
		c.Data["Breeds"] = nil
	} else {
		// Otherwise, assign the breed data to the template data
		c.Data["Breeds"] = breeds
		c.Data["Name"] = breeds[0].ID
		c.Data["Name"] = breeds[0].Name
		c.Data["Desc"] = breeds[0].Desc
		c.Data["Origin"] = breeds[0].Origin
		c.Data["Wiki"] = breeds[0].Wiki
	}

	// Handle the image URL
	if catImageURL == "" {
		// If the image URL is empty (i.e., an error occurred)
		c.Data["CatImageURL"] = "Error fetching image"
	} else {
		// Otherwise, assign the image URL to the template data
		c.Data["CatImageURL"] = catImageURL
	}

	// Render the HTML template with the image and breed data
	c.TplName = "index.html"
}

func (c *MainController) FetchCatImages() {
	// Retrieve breed_id from query parameters
	breedID := c.GetString("breed_id")
	if breedID == "" {
		breedID = "abys"
	}
	fmt.Println("Id:", breedID)

	// Create search URL with limit=8
	url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_ids=%s&limit=8", breedID)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
		c.ServeJSON()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Define a struct for parsing the image data
	var images []struct {
		URL string `json:"url"`
	}

	// Unmarshal the JSON response into the images struct
	err = json.Unmarshal(body, &images)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Return the images as JSON response
	if len(images) > 0 {
		c.Data["json"] = images
	} else {
		c.Data["json"] = map[string]interface{}{
			"error": "No images found for the selected breed",
		}
	}
	c.ServeJSON()
}
