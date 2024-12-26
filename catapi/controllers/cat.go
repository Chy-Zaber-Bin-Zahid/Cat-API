package controllers

import (
	"bytes"
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
	ID  string `json:"id"`
    URL string `json:"url"`
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
    breedChan := make(chan []Breed)
    imageChan := make(chan *CatImage) // Channel to send image data (ID and URL)

    // Fetch breed data in a goroutine
    go func() {
        apiURL := "https://api.thecatapi.com/v1/breeds"
        resp, err := http.Get(apiURL)
        if err != nil {
            fmt.Println("Error fetching breed data:", err)
            breedChan <- nil
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error reading response body:", err)
            breedChan <- nil
            return
        }

        var breeds []Breed
        err = json.Unmarshal(body, &breeds)
        if err != nil {
            fmt.Println("Error parsing JSON:", err)
            breedChan <- nil
            return
        }

        breedChan <- breeds
    }()

    // Fetch cat image data (ID and URL) in a goroutine
    go func() {
        apiURL := "https://api.thecatapi.com/v1/images/search"
        resp, err := http.Get(apiURL)
        if err != nil {
            fmt.Println("Error fetching cat image:", err)
            imageChan <- nil
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("Error reading response body:", err)
            imageChan <- nil
            return
        }

        var catImages []CatImage
        err = json.Unmarshal(body, &catImages)
        if err != nil || len(catImages) == 0 {
            fmt.Println("Error parsing JSON:", err)
            imageChan <- nil
            return
        }

        imageData := &CatImage{
            ID:  catImages[0].ID,
            URL: catImages[0].URL,
        }

        imageChan <- imageData
    }()

    // Get the breed and image data from the channels
    breeds := <-breedChan
    catImage := <-imageChan

    // Send breed and image data to the template
    if breeds == nil {
        c.Data["Breeds"] = nil
    } else {
        c.Data["Breeds"] = breeds
    }

    if catImage == nil {
        c.Data["CatImageID"] = "Error fetching image ID"
        c.Data["CatImageURL"] = "Error fetching image URL"
    } else {
        c.Data["CatImageID"] = catImage.ID
        c.Data["CatImageURL"] = catImage.URL
    }

    // Render the HTML template
    c.TplName = "index.html"
}

type Favorite struct {
    ImageID   string `json:"image_id"`
}

func (c *MainController) AddToFavourites() {
    var favoriteReq struct {
        ImageID string `json:"image_id"`
        SubID   string `json:"sub_id"`
    }

    // Validate incoming request body
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &favoriteReq); err != nil {
        c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
        c.Data["json"] = map[string]string{"error": "Invalid request payload"}
        c.ServeJSON()
        return
    }

    // Proceed with the external API request to add to favorites
    favChan := make(chan Favorite)
    errChan := make(chan error)

    go func() {
        apiKey, err := beego.AppConfig.String("apiKey")
        if err != nil || apiKey == "" {
            errChan <- fmt.Errorf("API key is missing or invalid")
            return
        }

        // Prepare the payload for the POST request
        body, err := json.Marshal(map[string]string{
            "image_id": favoriteReq.ImageID,
            "sub_id":   favoriteReq.SubID,
        })
        if err != nil {
            errChan <- err
            return
        }

        req, err := http.NewRequest("POST", "https://api.thecatapi.com/v1/favourites", bytes.NewBuffer(body))
        if err != nil {
            errChan <- err
            return
        }

        req.Header.Add("x-api-key", apiKey)
        req.Header.Add("Content-Type", "application/json")

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errChan <- err
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
            errChan <- fmt.Errorf("failed to add to favourites: %s", resp.Status)
            return
        }

        var favoriteResp Favorite
        if err := json.NewDecoder(resp.Body).Decode(&favoriteResp); err != nil {
            errChan <- err
            return
        }

        favChan <- favoriteResp
    }()

    // Handle the response from the goroutine
    select {
    case favorite := <-favChan:
        c.Data["json"] = favorite
        c.ServeJSON()
    case err := <-errChan:
        c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
    }
}



func (c *MainController) FetchCatImages() {
	// Retrieve breed_id from query parameters
	breedID := c.GetString("breed_id")
	if breedID == "" {
		breedID = "abys"
	}
	fmt.Println("Id:", breedID)

	// Fetch images from the first API
	firstAPIURL := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_ids=%s&limit=8", breedID)
	images, err := fetchImages(firstAPIURL)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Fetch breed information from the second API
	secondAPIURL := fmt.Sprintf("https://api.thecatapi.com/v1/breeds/%s", breedID)
	breedInfo, err := fetchBreedInfo(secondAPIURL)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Combine results into a single response
	response := map[string]interface{}{
		"images":    images,
		"breedInfo": breedInfo,
	}

	// Return the combined response as JSON
	c.Data["json"] = response
	c.ServeJSON()
}

// fetchImages fetches images from the API and parses the JSON response
func fetchImages(url string) ([]struct {
	URL string `json:"url"`
}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var images []struct {
		URL string `json:"url"`
	}
	err = json.Unmarshal(body, &images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// fetchBreedInfo fetches breed-specific information from the API and parses the JSON response
func fetchBreedInfo(url string) (struct {
	Name        string `json:"name"`
	Origin      string `json:"origin"`
	Description string `json:"description"`
	WikipediaURL string `json:"wikipedia_url"`
}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return struct {
			Name        string `json:"name"`
			Origin      string `json:"origin"`
			Description string `json:"description"`
			WikipediaURL string `json:"wikipedia_url"`
		}{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return struct {
			Name        string `json:"name"`
			Origin      string `json:"origin"`
			Description string `json:"description"`
			WikipediaURL string `json:"wikipedia_url"`
		}{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return struct {
			Name        string `json:"name"`
			Origin      string `json:"origin"`
			Description string `json:"description"`
			WikipediaURL string `json:"wikipedia_url"`
		}{}, err
	}

	var breedInfo struct {
		Name        string `json:"name"`
		Origin      string `json:"origin"`
		Description string `json:"description"`
		WikipediaURL string `json:"wikipedia_url"`
	}
	err = json.Unmarshal(body, &breedInfo)
	if err != nil {
		return breedInfo, err
	}

	return breedInfo, nil
}


