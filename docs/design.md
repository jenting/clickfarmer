# Click Farmer Design Proposal

## Overview

The click farmer will have two main components. The primary component will be
referred to as the "database", and it keeps track of a global count of red,
green, and blue clicks. It also hosts the GRPC server, which web servers can use
to get and set values for red, green, and blue clicks.

The other component, referred to as the "webserver", consists of three important
parts. Firstly, it serves a simple webpage with red, green, and blue buttons,
which allows us to farm clicks from visitors. Secondly, it hosts an HTTP API
which the webpage reaches out to in order to update click counts. Finally, it
runs a GRPC client, which allows the client to get and update the global count
of the central server.

We plan to run one database at a time, but we're in the business of farming clicks
from all over the world. So the same limitation does not apply to the webserver.
We will run lots of webservers against the database. The webserver will maintain
a local cache of click counts and periodically sync these values with the
database, to precisely collect clicks with low load.

## Server

The RPC server functionality is defined as follows:

```
service ClickFarmer {
    rpc GetClicks(GetClicksRequest) returns (GetClicksResponse) {}
    rpc SetClicks(SetClicksRequest) returns (SetClicksResponse) {}
}

message ClickCounts {
    int64 red = 1;
    int64 green = 2;
    int64 blue = 3;
}

message GetClicksRequest {}

message GetClicksResponse {
    ClickCounts clickCounts = 1;
}

message SetClicksRequest {
    ClickCounts clickCounts = 1;
}

message SetClicksResponse {}
```

The `GetClicks` RPC call will return the server's current counts for red, green,
and blue clicks.

The `SetClicks` RPC call will set the server's current counts for red, green, and
blue clicks.

## Client

The HTTP API of the client functions as follows:

* `GET /api/clicks` - The client will return the values of the clicks in the
    local cache as a JSON object in an http response.
* `GET /api/clicks/red` - The client will return the number of red clicks in the
    local cache as an http response.
* `GET /api/clicks/green` - The client will return the number of green clicks in
    the local cache as an http response.
* `GET /api/clicks/blue` - The client will return the number of blue clicks in
    the local cache as an http response.
* `PUT /api/clicks/red` - The client will increase red clicks in the local cache
    by 1.
* `PUT /api/clicks/green` - The client will increase green clicks in the local
    cache by 1.
* `PUT /api/clicks/blue` - The client will increase blue clicks in the local
    cache by 1.

The local cache should refresh on a configurable interval. The steps for a
refresh are as follows:

* Call `SetClicks` on the RPC server to update the values.
* Call `GetClicks` on the RPC server to get updated values from the server.

