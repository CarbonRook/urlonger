# carbonrook/urlonger

Package `carbonrook/urlonger` implements an HTTP redirect chain resolver for identifying the final destination of a shortened URL.

Current implementation has some limitations/design decisions:
* Redirects are resolved using the `HEAD` HTTP request
* JavaScript redirects are not evaluated
* For simplicity of output only filtered headers are returned

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/carbonrook/urlonger
```

# CLI Examples

```text
./bin/urlonger -help
Usage of ./bin/urlonger:
  -o string
    	output format; text, json (default "json")
  -url string
    	url to resolve
  -v	enable verbose logging
```

Standard usage:

```bash
urlonger -url https://bit.ly/3Bg19uM
```

Verbose logging:

```bash
urlonger -url https://bit.ly/3Bg19uM -v
```

## Examples

Output the redirect chain in a JSON format using the Resolve function:

```go
func main() {

	redirs, err := urlonger.Resolve("https://bit.ly/3Bg19uM", []string{"server"})
	if err != nil {
		log.Fatal(err)
	}

    redirsJson, err := json.MarshalIndent(redirs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(redirsJson))
}
```

Default output is a JSON array of structs in the order that the request and redirect occurs. A response which does not include a `destination` header is considered to be the final URL in the chain of redirects:

```json
[
  {
    "url": "https://bit.ly/3Bg19uM",
    "status_code": 301,
    "destination": "https://www.usatoday.com/story/travel/2022/02/10/amtrak-deal-valentines-offer-sale/6741296001/",
    "response_headers": {
      "server": "nginx",
      "via": "1.1 google"
    }
  },
  {
    "url": "https://www.usatoday.com/story/travel/2022/02/10/amtrak-deal-valentines-offer-sale/6741296001/",
    "status_code": 302,
    "destination": "https://eu.usatoday.com/story/travel/2022/02/10/amtrak-deal-valentines-offer-sale/6741296001/",
    "response_headers": {}
  },
  {
    "url": "https://eu.usatoday.com/story/travel/2022/02/10/amtrak-deal-valentines-offer-sale/6741296001/",
    "status_code": 200,
    "destination": "",
    "response_headers": {
      "via": "1.1 varnish, 1.1 varnish"
    }
  }
]
```