# Web Crawler

## Scope
Given a starting URL, the crawler should visit each URL it finds on the same domain. It should print each URL visited, and a lists of links found on each page. The crawler will be limited to one subdomain, but not follow external links or subdomains that don't match the provided domain. 

For example the given domain is `https://monzo.com`, the page has URLs suchs as `facebook.com` or `community.monzo.com`, the crawler will not follow these links

## Decisions made
The solution is written with an entrypoint being a CLI but the main functionality has been seprataed out in order to make it useable in other entrypoints such as an API.

## Architecture / Diagrams


## Prerequisites
1.  Golang 1.24 / Docker 


## Running the CLI
The CLI has the following options:

### Running locally
To run locally you will need to execute the following commands:
```bash
make build
TODO fill in example
```

### Running via Docker
Run the following:
```bash
make docker
TODO fill in example
```


## Note / Further Improvements
TODO add info here


### Further improvements
1.  Websocket API could be added to the functionality to push updates to a webpage as new page is found with the URL.