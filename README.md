HttpServer project is a web server that accepts "GET" requests, looks for malformed URLs and replies to sender with a response about whether the URL was malformed or not.

Use build.sh to build the Docker image.
Use run.sh to run the web server.

Tested it with:
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/hello
-> expected result 200 OK "Hello"

and
curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/hellob
-> expected result 404
