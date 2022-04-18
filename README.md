# Click Farmer

Make sure to read the "Additional Notes" section at the bottom for some important
information.

## The Problem

The click farming business has really taken off recently.

The original click farming engineers started this project to count the number of
clicks around the world. However, they left to start another business.

They created a design document that describes how the system works, but they never
got around to finishing. We're not quite sure whether everything is correct either.

They say the system they left for you consists of a core database service they
wrote. In addition, they wrote a webserver. They expect to run many webservers
against one database concurrently, to handle the clickfarming load.

In the database folder you'll find the RPC server implementation that handles
database requests.

In the webserver folder you'll find the HTTP handlers and code which talks to
the database using protobuf + grpc defined in the pb folder.

In the docs folder you'll find the design document that the click farming
engineers wrote.

The goal is to get the product working to farm upvotes from around the world.

## Building and Running

To build the binary:
```
go build -o clickfarmer -race -v .
```

To run the database:
```
./clickfarmer database
```
You can specify an `--rpc-addr <addr>` flag to set the RPC server address the
database listens for requests on to something other than ":8080".

To run a webserver:
```
./clickfarmer webserver
```
You can specify an `--rpc-addr <addr>` flag to set the address to connect to
the database on (it defaults to "localhost:8080").

You can specify an `--http-addr <addr>` (or `-a`) flag to set the address to
serve the website on. This is used to run multiple web servers, like:

```
./clickfarmer webserver -a :3001
./clickfarmer webserver -a :3002
./clickfarmer webserver -a :3003
```

Note that you can also run the database and webserver without the `build` step
by running the `main.go` file directly:

```
go run main.go database
go run main.go webserver -a :3002
```

To access the webserver's frontend, go to http://localhost:3000 in a browser (if
you used `--http-addr` or `-a` to specify a different one, use that port
instead of `3000`).

## Generating Protobufs

You will need to install the protobuf compiler:

https://grpc.io/docs/protoc-installation/

As well as some additional dependencies:

```
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

To re-generate the protobuf code after changing `pb/clicktracking.proto`, run:
```
go generate ./pb/...
```

## Additional Notes

In order to get full credit you must make sure the following endpoints work
on the webserver API:
```
GET /api/clicks/red -> returns a single integer representing total red button clicks (`0`)
GET /api/clicks/green -> returns a single integer representing total green button clicks
GET /api/clicks/blue -> returns a single integer representing total blue button clicks
GET /api/clicks -> returns JSON like `{"redClicks": 0, "greenClicks": 0, "blueClicks": 0}`
PUT /api/clicks/red -> called when the red button was clicked once. no return data
PUT /api/clicks/green -> called when the green button was clicked once. no return data
PUT /api/clicks/blue -> called when the blue button was clicked once. no return data
```

You must match this behavior for all of these endpoints - we use them to test
your solution.

You also must keep the same command line functionality. Our testing software
expects to be able to start a webserver like
`./clickfarmer webserver -a :port` and expects to start a database like
`./clickfarmer database`. It expects that each webserver syncs with the database
every 1 second.

You must keep the webserver/database architecture. Our testing software expects
to run multiple webservers against a database. However, you can do whatever
you want in terms of the webserver and database's internal logic, as well as the
messages and RPC functions they use to communicate.

There are two components to your homework evaluation:
* Backend - we will run automated tests to determine whether your API is able to correctly handle requests, including when multiple APIs are running against the same database
* Frontend - we will look at the html page served by your web server and assign points based on responsiveness to different screen sizes, and design decisions such as color and element spacing/positioning

