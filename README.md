
# GoQueryBlogs

GoQueryBlogs is a REST API written in Go. The API fetches blog posts from the URL provided by the user. The API response will be a JSON object with a list of all blog posts that have at least one tag belongs to the tags parameter specified in the query string. 

## Query parameters
The `tags` parameter is required in the query string. The `sortBy` and `direction` parameters are optional. The acceptable `sortBy` fields are: id, reads, likes, popularity. The acceptable `direction` fields are: asc, desc. By default, the API will sort the JSON blog posts by id in ascending order. The API makes concurrent requests to the URL. QueryBlogs has also implemented a server side cache to reduce the number of calls made.

## Getting Started

To build Go executable, run:
```sh
go build
```
To run the API, type the following command:
```sh
go run .
```
To execute unit test and integration test, run:
```sh
go test
```
## Running the API locally
GoQueryBlogs use port 8080 to run the application. To test if the server is running correctly, type the following command:
```sh
curl http://localhost:8080/api/ping
```
To fetch blog posts, type the following command:
```sh
curl -i -X GET 'http://localhost:8080/api/posts?tags=tech,science&sortBy=likes&direction=desc'
```

## Docker
You can also deploy the API with docker.
To build a docker image, run:
```sh
docker build . -t go-dock
```
To run an image, type the following command:
```sh
docker run -p 8080:8080 go-dock
```
If you want to bind to another port, you can run the above command with `-p $HOST_PORT:8080`





