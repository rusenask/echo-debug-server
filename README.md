# HTTP debug echo server

Send requests at it and it will disect them back. 

Start the server:

```
go run main.go
```

Then, any incoming requests will be reflected back with the details:

```
curl -X POST http://localhost:8888\?some\=query --data '{"json": "data"}' 
{
  "ServerInfo": {
    "Port": ":8888"
  },
  "method": "POST",
  "path": "/",
  "raw_query": "some=query",
  "header": {
    "Accept": [
      "*/*"
    ],
    "Content-Length": [
      "16"
    ],
    "Content-Type": [
      "application/x-www-form-urlencoded"
    ],
    "User-Agent": [
      "curl/7.68.0"
    ]
  },
  "body": "{\"json\": \"data\"}",
  "form_values": null
}
```
