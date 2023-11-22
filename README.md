# Overview:
This program implements a simple HTTP server using the Go programming language. The server listens on a specified port and serves files from a specified directory. There are additional functionalities, like handling requests for specific routes or files, logging requests, and handling errors.

## Features and Functionalities:
### Request Handling:
- If the request doesn't follow the basic structure of an HTTP request, a 400 (Bad Request) response is sent.
- Requests to the root ("/") are automatically redirected to "/index.html".
- Any request starting with "/private" is treated as a forbidden resource, and a 403
  (Forbidden) response is sent.
- If the requested file exists in the specified directory, it is served to the client. Otherwise, a
  404 (Not Found) response is sent.
- For other errors, a 500 (Internal Server Error) is sent.

### Dispatcher Mechanism:
The server continuously listens for incoming client connections using the ln.Accept() method. When a client connects to the server, it establishes a new connection represented by the conn object.

For each accepted connection, the server doesn't process the request in the main execution thread. Instead, it spawns a new lightweight thread, known as a "goroutine", to handle the client's request. This is done using the go keyword followed by the processClientRequest(conn) function. The use of goroutines allows the server to handle multiple requests simultaneously without waiting for one request to complete before moving to the next.

### Content Type Handling: 
The server can serve various file types, such as .html, .txt, .jpg, and .gif.

### Custom Port and Document Root: 
Through the use of command-line flags, the server can be started on a custom port (-port) and can serve files from a specified directory (-document_root).

### Client Request Handling:
When a client connects to the server and sends an HTTP request, the processClientRequest function reads the request, extracts the necessary headers, and delegates further handling to the handleRequest function.

### Timezone Handling: 
The program is set to recognize the Pacific timezone (America/Los Angeles) for logging purposes.

### Logging:
Every response from the server is logged to the console with a timestamp in the Pacific time zone, the HTTP status code, and the request URI.

# Usage:
```
go run server.go -port=8888 -document_root=./webserver_files
```
This would start the server on port 8888 and serve files from the "webserver_files" directory.

### Screenshots attached in folder
- Url: http://localhost:8888/ SCU Home Page; Status Code: 200
- Url: http://localhost:8888/private; Private Url for which access is denied; Status Code: 403 
- Url: http://localhost:8888/indexerror.html; url exists server unable to serve. Status Code: 500 (To work do "chmod 000 indexerror.html")
- Supporting txt; requirements.txt http://localhost:8888/requirements.txt; Status Code: 200
- SCU GIF; url: http://localhost:8888/scu.gif; Status Code: 200
- SCU JPG: http://localhost:8888/scu.jpg; Status Code: 200
- Status Code: 400 Using Terminal: echo -ne "GET" | nc localhost 8888
