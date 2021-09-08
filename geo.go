package geo

import (
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	geoip2 "github.com/oschwald/geoip2-golang"
)

// Error defines the error in returning Geographical information
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response defines the structure of Geographical information
type Response struct {
	IPAddress     string  `json:"IPAddress"`
	CityName      string  `json:"CityName"`
	StateCode     string  `json:"StateCode"`
	CountryCode   string  `json:"CountryCode"`
	ContinentCode string  `json:"ContinentCode"`
	TimeZone      string  `json:"TimeZone"`
	ZipCode       string  `json:"ZipCode"`
	Latitude      float64 `json:"Latitude"`
	Longitude     float64 `json:"Longitude"`
	Language      string  `json:"Language"`
	Error         Error   `json:"error"`
}

// getErrorResponse returns the error if something goes wrong
func getErrorResponse(errResponse string) *Response {
	err := Error{
		Code:    400,
		Message: errResponse,
	}
	return &Response{
		Error: err,
	}
}

// getResponse Maps the record from Maxmind in appropriate format
func getResponse(ipAddress string, language string, db *geoip2.Reader) *Response {
	ip := net.ParseIP(ipAddress)
	record, err := db.City(ip)
	if err != nil {
		return getErrorResponse("Could not get Geo information")
	}

	var stateCode string
	if len(record.Subdivisions) > 0 {
		stateCode = record.Subdivisions[0].Names[language]
	}

	return &Response{
		IPAddress:     ipAddress,
		CityName:      record.City.Names[language],
		StateCode:     stateCode,
		CountryCode:   record.Country.IsoCode,
		ContinentCode: record.Continent.Code,
		TimeZone:      record.Location.TimeZone,
		ZipCode:       record.Postal.Code,
		Latitude:      record.Location.Latitude,
		Longitude:     record.Location.Longitude,
		Language:      language,
	}
}

// getLanguage returns the language of the user from the header
func getLanguage(c *gin.Context) string {
	acptLang := strings.TrimSpace(c.Request.Header.Get("ACCEPT-LANGUAGE"))
	geoIPSupported := []string{"fr", "de", "ja", "ru", "es", "pt-BR", "zh-CN", "en"}
	for _, lang := range geoIPSupported {
		if strings.Contains(acptLang, lang) {
			return lang
		}
	}
	return "en"
}

// getClientIP returns the IP Address of the user from the headers
func getClientIP(c *gin.Context) (string, error) {
	xForwardedFor := strings.TrimSpace(c.Request.Header.Get("X-FORWARDED-FOR"))
	remoteAddr := strings.TrimSpace(c.Request.Header.Get("REMOTE-ADDR"))
	clientIP := strings.TrimSpace(c.Request.Header.Get("CLIENT-IP"))

	ipAddr := ""
	if len(xForwardedFor) != 0 {
		ipAddr = xForwardedFor
	} else if len(remoteAddr) != 0 {
		ipAddr = remoteAddr
	} else if len(clientIP) != 0 {
		ipAddr = clientIP
	}
	if len(ipAddr) != 0 {
		ip := net.ParseIP(ipAddr)
		if ip == nil || ip.IsLoopback() {
			return "", errors.New("Invalid IP or Loopback IP address")
		}
		ip = ip.To4()
		if ip == nil {
			return "", errors.New("Could not get IPv4 address")
		}
		return ipAddr, nil
	}
	return "", errors.New("Could not get client IP address")
}

// setContext sets the geographical information in Gin context
func setContext(c *gin.Context, db *geoip2.Reader) {
	start := time.Now()
	ipAddress, err := getClientIP(c)
	if err == nil {
		language := getLanguage(c)
		response := getResponse(ipAddress, language, db)
		c.Set("GeoResponse", response)
	} else {
		response := getErrorResponse(err.Error())
		c.Set("GeoResponse", response)
	}
	duration := time.Now().Sub(start)
}

// getDB returns the database handle
func getDB(dbPath string) (*geoip2.Reader, *Response) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, getErrorResponse("Maxmind DB not found")
	}
	return db, nil
}

// Middleware sets the Geographical information
// about the user in the Gin context
func Middleware(dbPath string) gin.HandlerFunc {
	db, dbErr := getDB(dbPath)
	return func(c *gin.Context) {
		if dbErr == nil {
			setContext(c, db)
		}
		c.Next()
	}
}

// Default returns the handler that sets the
// geographical information about the user
func Default(dbPath string) gin.HandlerFunc {
	return Middleware(dbPath)
}
