# PUBSUB

This publish/subscribe server uses [Server Sent Events](sses).

The server is simple enough to run:
```
%> go run main.go
```
You can subscribe to messages via HTTP:
```
%> curl -N localhost:9999/subscribe
```
You can publish a message via HTTP as well:
```
%> curl -d "this is my data" localhost:9999/publish
```
Tests:
```
%> go test -race
%> errcheck .
```
[sses]: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events