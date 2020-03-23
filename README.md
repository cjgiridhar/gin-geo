# gin-geo

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

[Gin](https://github.com/gin-gonic/gin) middleware that sets the geographical information about the user in the Gin context

## Project Context and Features

Detecting the location of users visiting your website is useful for a variety of reasons. 

- You might want to display different content based on different languages for people from different countries OR 
- Display targeted information to visitors from different locations OR 
- You may have setup database shards and want user signups to go to appropriate database based on their geography

Whatever your reasons, this middleware comes to your rescue!

## Requirements

gin-geo uses the following [Go](https://golang.org/) packages as
dependencies:

- github.com/gin-gonic/gin

## Installation

Step 1: Assuming you've installed Go and Gin, run this command to get the middleware:  
```
$ go get github.com/cjgiridhar/gin-geo
```

Step 2: Download latest GeoLite2 City Database from Maxmind.
```
$ wget http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
$ gunzip GeoLite2-City.mmdb.gz  
```
OR
Copy it from the gin-geo repository
```
cp github.com/cjgiridhar/gin-geo/example/GeoLite2-City.mmdb .
```

Step 3: Make sure file path for mmdb file is at correct location.
Please place the file from where your program runs.

If you look at the [gin-geo](https://github.com/cjgiridhar/gin-geo/tree/master/example) sample code, 
as ```example.go``` runs from ```example``` folder, we have placed the database at that location.

## Usage

main.go
```go
package main

import (
	"fmt"
	"net/http"

	geo "github.com/cjgiridhar/gin-geo"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(geo.Default())
	r.GET("/geo", func(c *gin.Context) {
		geoResponse, ok := c.Get("GeoResponse")
		if ok {
			fmt.Println(geoResponse)
			c.JSON(http.StatusOK, gin.H{
				"geo": geoResponse,
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Could not get geographical information",
			})
		}
	})
	r.Run()
}
```

## Run the demo

You need a request from the internet to see the middleware in action.
For this you will need a public URL for exposing your local web server. 

Follow the steps:
- Create a ```main.go``` file and copy the contents as given in the Usage section.
- Make sure GeoLite2-City.mmdb is placed in the same path as main.go.
- Run the demo as: ```go run main.go``` (This will run a local web server on port 8080).
- Download ngrok and run ```./ngrok http 8080``` (This will generate public URL and expose local web server to the internet).
- Browse the public URL obtained from step 2, say https://public.ngrok.io/geo, from your browser.

## Sample Output

```
{
	"geo": {
		"IPAddress":"151.236.26.140",
		"CityName":"Zurich",
		"StateCode":"Zurich",
		"CountryCode":"CH",
		"ContinentCode":"EU",
		"TimeZone":"Europe/Zurich",
		"ZipCode":"8048",
		"Latitude":47.3667,
		"Longitude":8.55,
		"Language":"de",
	}
}
```

## Server Logs

```
2020/03/18 09:54:36 Geo: Middleware duration 391.525µs
&{151.236.26.140 Zurich Zurich CH EU Europe/Zurich 8048 47.3667 8.55 de {0 }}
[GIN] 2020/03/18 - 09:54:36 | 200 |     472.715µs |    103.4.18.143 | GET      "/geo"

```

## How can I use the middleware?

Based on the above sample output, you can use the CountryCode or ContinentCode information to make meaningful decisions for the client.

For instance, if it's a EU country, you can create the account for the client on a server hosted in EU region as per GDPR requirements.


## Thanks

This middleware uses Maxmind's GeoLite (https://www.maxmind.com/en/home) database to get geo information
and used HTTP headers to get the IP Address and Language for the user.

## Contact

Chetan Giridhar, cjgiridhar@gmail.com