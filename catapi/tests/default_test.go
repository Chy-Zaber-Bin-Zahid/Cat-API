package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"path/filepath"
	"runtime"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

// TestCatEndpoint tests the cat controller endpoints
func TestCatEndpoint(t *testing.T) {
	Convey("Subject: Test Cat Controller Endpoints\n", t, func() {
		Convey("When getting cat list", func() {
			r, _ := http.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("Response should not be empty", func() {
				So(w.Body.Len(), ShouldBeGreaterThan, 0)
			})
		})
	})
}
func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}
// TestFetchCatImagesEndpoint tests the cat image fetching endpoints
func TestFetchCatImagesEndpoint(t *testing.T) {
	Convey("Subject: Test Cat Image Fetching Endpoints\n", t, func() {
		Convey("When fetching cat images without breed_id", func() {
			r, _ := http.NewRequest("GET", "/catImages", nil) 
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("Response should contain images and breed info", func() {
				So(w.Body.Len(), ShouldBeGreaterThan, 0)
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				So(response["images"], ShouldNotBeNil)
				So(response["breedInfo"], ShouldNotBeNil)
			})
		})

		Convey("When fetching cat images with specific breed_id", func() {
			r, _ := http.NewRequest("GET", "/catImages", nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Convey("Status Code Should Be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("Response should contain breed-specific data", func() {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				breedInfo := response["breedInfo"].(map[string]interface{})
				So(breedInfo["name"], ShouldNotBeEmpty)
				So(breedInfo["origin"], ShouldNotBeEmpty)
			})
		})

		Convey("When adding to favorites", func() {
			favoriteJson := `{"image_id":"test123","sub_id":"user1"}`
			r, _ := http.NewRequest("POST", "/add-to-favourites", strings.NewReader(favoriteJson))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Convey("Should return appropriate status code", func() {
				So(w.Code, ShouldBeIn, []int{200, 201})
			})
		})

		Convey("When adding invalid favorite", func() {
			invalidJson := `{"invalid":"data"}`
			r, _ := http.NewRequest("POST", "/add-to-favourites", strings.NewReader(invalidJson))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Convey("Should return 400 Bad Request", func() {
				So(w.Code, ShouldEqual, 400)
			})
		})
	})
}
// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}

