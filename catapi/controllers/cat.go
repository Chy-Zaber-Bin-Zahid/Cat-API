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
	Image  struct {
		URL string `json:"url"`
	} `json:"image"`
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
