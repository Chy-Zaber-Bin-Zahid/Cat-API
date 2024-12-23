package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

// Struct to parse the JSON response from the API
type CatImage struct {
	URL string `json:"url"`
}

func (c *MainController) Get() {
	// Channel to get the image URL
	imageChan := make(chan string)

	// Start the goroutine to fetch the image URL from the external API
	go func() {
		// Format the API URL with the key
		apiURL := "https://api.thecatapi.com/v1/images/search"

		// Make the HTTP request
		resp, err := http.Get(apiURL)
		if err != nil {
			// If there's an error with the HTTP request
			fmt.Println("Error fetching cat image:", err)
			imageChan <- ""
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// If there's an error reading the response body
			fmt.Println("Error reading response body:", err)
			imageChan <- ""
			return
		}

		// Unmarshal the JSON response
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

	// Fetch the image URL from the channel
	catImageURL := <-imageChan
	if catImageURL == "" {
		// If the image URL is empty (i.e., an error occurred)
		c.Data["CatImageURL"] = "Error fetching image"
	} else {
		// Otherwise, assign the image URL to the template data
		c.Data["CatImageURL"] = catImageURL
	}

	// Render the HTML template with the image URL
	c.TplName = "index.html"
}
