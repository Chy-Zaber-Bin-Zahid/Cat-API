package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

type CatImage struct {
	URL string `json:"url"`
}

func (c *MainController) Get() {
	// Check if the request expects a JSON response
	acceptHeader := c.Ctx.Request.Header.Get("Accept")
	if acceptHeader == "application/json" {
		// Return JSON response with images
		apiURL := "https://api.thecatapi.com/v1/images/search?limit=10&breed_ids=beng&api_key=live_Ii20w7Wt785t9kCsxDQYAMTIIL7epsK1IaGiHL3hxWw0ou2AfkvZ3FAMxJ4NEc0Z"
		ch := make(chan []string)

		var wg sync.WaitGroup
		wg.Add(1)

		// Fetching images concurrently
		go func() {
			defer wg.Done()
			response, err := http.Get(apiURL)
			if err != nil {
				fmt.Println("Error fetching data:", err)
				ch <- nil
				return
			}
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				ch <- nil
				return
			}

			var images []CatImage
			if err := json.Unmarshal(body, &images); err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				ch <- nil
				return
			}

			var imageURLs []string
			for _, img := range images {
				imageURLs = append(imageURLs, img.URL)
			}
			ch <- imageURLs
		}()

		go func() {
			wg.Wait()
			close(ch)
		}()

		// Receiving image URLs
		imageURLs := <-ch
		if imageURLs == nil {
			imageURLs = []string{"https://placekitten.com/500/300"} // Fallback image
		}

		// Send the image URLs as JSON response
		c.Ctx.Output.SetStatus(200)
		c.Ctx.Output.Header("Content-Type", "application/json")
		c.Ctx.Output.JSON(map[string]interface{}{
			"images": imageURLs,
		}, false, false)

	} else {
		// Default behavior for the initial page load (HTML render)
		apiURL := "https://api.thecatapi.com/v1/images/search?limit=10&breed_ids=beng&api_key=live_Ii20w7Wt785t9kCsxDQYAMTIIL7epsK1IaGiHL3hxWw0ou2AfkvZ3FAMxJ4NEc0Z"
		ch := make(chan []string)

		var wg sync.WaitGroup
		wg.Add(1)

		// Fetching images concurrently
		go func() {
			defer wg.Done()
			response, err := http.Get(apiURL)
			if err != nil {
				fmt.Println("Error fetching data:", err)
				ch <- nil
				return
			}
			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error reading response:", err)
				ch <- nil
				return
			}

			var images []CatImage
			if err := json.Unmarshal(body, &images); err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				ch <- nil
				return
			}

			var imageURLs []string
			for _, img := range images {
				imageURLs = append(imageURLs, img.URL)
			}
			ch <- imageURLs
		}()

		go func() {
			wg.Wait()
			close(ch)
		}()

		// Receiving image URLs
		imageURLs := <-ch
		if imageURLs == nil {
			imageURLs = []string{"https://placekitten.com/500/300"} // Fallback image
		}

		// Pass data to template
		c.Data["Images"] = imageURLs
		c.TplName = "index.html"
	}
}
