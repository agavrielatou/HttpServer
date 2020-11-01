HttpServer project is an HTTP server that accepts "GET" requests. For every request, it looks at a Postgres database to check if the given URL is included in the malformed ones and replies to sender with a response about whether the URL is valid or not. Server works with a JSON configuration file (config.json). Also long URLs are encoded using base64.

Use build.sh to build the Docker image.
Use run.sh to run the web server.

Dockerfile for the Postgres database is also included. Run postgres/build.sh to build the image and postgres/run.sh to run it. It adds a table with 5 malformed urls for the purpose of this demo.

Tested it using Postman.
Was also able to test it with AWS RDS Postgres database.

Followed step by step additions and commits.
My first project in Go!


