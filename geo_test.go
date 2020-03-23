package geo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	fp = "example/GeoLite2-City.mmdb"
)

// Tests the Response structure
func TestResponse(t *testing.T) {
	response := &Response{
		IPAddress: "1.1.1.1",
	}
	if response == nil {
		t.Errorf("Response structure failed, got %s, want %s",
			response.IPAddress, "1.1.1.1")
	}
}

// Tests the Error interface
func TestErrorResponse(t *testing.T) {
	response := getErrorResponse("Error string")
	if response == nil {
		t.Errorf("Response structure failed, got %s, want %s",
			response.Error.Message, "Error string")
	}
}

// Tests if the Client IP is detected correctly by middleware
// from X-FORWARDED-FOR header
func TestClientIPForwardedHeader(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	c.Request.Header.Set("X-FORWARDED-FOR", "1.1.1.1")
	ipAddr, _ := getClientIP(c)
	if ipAddr != "1.1.1.1" {
		t.Errorf("Getting Client IP is incorrect, got %s, want %s.",
			ipAddr, "1.1.1.1")
	}
}

// Tests if the Client IP is detected correctly by middleware
// from REMOTE-ADDR header
func TestClientIPRemoteAddrHeader(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	c.Request.Header.Set("REMOTE-ADDR", "1.1.1.1")
	ipAddr, _ := getClientIP(c)
	if ipAddr != "1.1.1.1" {
		t.Errorf("Getting Client IP is incorrect, got %s, want %s.",
			ipAddr, "1.1.1.1")
	}
}

// Tests if the Client IP is detected correctly by middleware
// from CLIENT-IP header
func TestClientIPHeader(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	c.Request.Header.Set("CLIENT-IP", "1.1.1.1")
	ipAddr, _ := getClientIP(c)
	if ipAddr != "1.1.1.1" {
		t.Errorf("Getting Client IP is incorrect, got %s, want %s.",
			ipAddr, "1.1.1.1")
	}
}

// Tests if the langauge is correctly detected by middlewire
func TestLanguage(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	c.Request.Header.Set("ACCEPT-LANGUAGE", "fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5")
	langauge := getLanguage(c)
	if langauge != "fr" {
		t.Errorf("Getting Client IP is incorrect, got %s, want %s.",
			langauge, "fr")
	}
}

// Tests if the invalid Client IP is handled correctly by middlewire
func TestClientNoIP(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	_, err := getClientIP(c)
	if err.Error() != "Could not get client IP address" {
		t.Errorf("Getting Client IP is incorrect, got %s, want %s.", err,
			"Could not get client IP address")
	}
}

// Tests if the middleware sets the context
// func TestMiddleware(t *testing.T) {
// 	buf := new(bytes.Buffer)
// 	client := &http.Client{}
// 	c, _ := gin.CreateTestContext(httptest.NewRecorder())
// 	c.Request, _ = http.NewRequest("GET", "https://baff60d4.ngrok.io/geo", buf)
// 	c.Request.Header.Set("X-FORWARDED-FOR", "123.123.123.123")
// 	resp, err := client.Do(c.Request)
// 	if err != nil {
// 		log.Fatal("Errored when sending request to the server")
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if err != nil {
// 		t.Errorf("Middleware failed, got %d, want %d.", resp.StatusCode, 200)
// 	}
// }

// Tests the middleware method, that sets the context, by creating Gin request
func TestGeo(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	c.Request.Header.Set("X-FORWARDED-FOR", "123.123.123.123")
	setContext(c, "example/GeoLite2-City.mmdb")
	geoResponse, _ := c.Get("GeoResponse")
	var r Response
	js, _ := json.Marshal(geoResponse)
	json.Unmarshal(js, &r)
	if r.CityName != "Beijing" {
		t.Errorf("Middleware worked, got %s, want %s.", r.CityName, "Beijing")
	}
}

// Tests the middleware method for invalid IP
func TestGeoInvalidIP(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	c.Request.Header.Set("X-FORWARDED-FOR", "")
	setContext(c, "example/GeoLite2-City.mmdb")
	geoResponse, _ := c.Get("GeoResponse")
	var r Response
	js, _ := json.Marshal(geoResponse)
	json.Unmarshal(js, &r)
	if r.Error.Code != 400 {
		t.Errorf("Middleware failed, got %d, want %d.", r.Error.Code, 400)
	}
}

// Tests the middleware method, when DB location is incorrect
func TestInvalidDBPath(t *testing.T) {
	buf := new(bytes.Buffer)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "/geo", buf)
	c.Request.Header.Set("X-FORWARDED-FOR", "123.123.123.123")
	setContext(c, "example/GeoLite2-City-NO.mmdb")
	geoResponse, _ := c.Get("GeoResponse")
	var r Response
	js, _ := json.Marshal(geoResponse)
	json.Unmarshal(js, &r)
	if r.Error.Message != "Maxmind DB not found" {
		t.Errorf("Middleware worked, got %s, want %s.",
			r.Error.Message, "Maxmind DB not found")
	}
}

// Tests if the Timezone is returned correctly by middleware
func TestTimeZone(t *testing.T) {
	clientIP := "121.72.165.118"
	mappedResponse := getResponse(clientIP, "en", fp)
	if mappedResponse.TimeZone != "Pacific/Auckland" {
		t.Errorf("Country Code is incorrect, got %s, want %s.",
			mappedResponse.TimeZone, "Pacific/Auckland")
	}
}

// Tests if the middleware detects if a country belongs to EU continent correctly
func TestEUCountry(t *testing.T) {
	clientIP := "104.238.171.182"
	mappedResponse := getResponse(clientIP, "en", fp)
	if mappedResponse.ContinentCode != "EU" {
		t.Errorf("EU Country check failed, got %s, want %s.",
			mappedResponse.ContinentCode, "EU")
	}
}

// Tests if the country code is returned correctly by middleware
func TestCountryCode(t *testing.T) {
	clientIP := "49.207.48.225"
	mappedResponse := getResponse(clientIP, "en", fp)
	if mappedResponse.CountryCode != "IN" {
		t.Errorf("Country Code is incorrect, got %s, want %s.",
			mappedResponse.CountryCode, "IN")
	}
}

// Tests if the middleware returns the correct city name
func TestCityName(t *testing.T) {
	clientIP := "123.123.123.123"
	mappedResponse := getResponse(clientIP, "en", fp)
	if mappedResponse.CityName != "Beijing" {
		t.Errorf("City is incorrect, got %s, want %s.",
			mappedResponse.CityName, "Beijing")
	}
}
