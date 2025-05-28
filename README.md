# httpfromtcp 🛜

httpfromtcp is a custom implementation of an HTTP/1.1 server built from the ground up using Go's `net` package. The primary goal of this project is to understand the fundamentals of the HTTP protocol by implementing its core features.

## Features

*   **HTTP/1.1 Compliance (Partial)**:
    *   Parses HTTP request lines, headers, and bodies.
    *   Constructs and sends HTTP responses including status lines, headers, and bodies.
*   **Request Routing**: Basic routing based on request path and method.
*   **Static File Serving**: Example endpoint (`/video`) to serve local video files.
*   **Proxying**: Example endpoint (`/httpbin/*`) that proxies requests to `httpbin.org`.
*   **Chunked Transfer Encoding**: Implemented for responses, particularly demonstrated in the proxy handler.
*   **Trailers**: Supports sending trailer headers after a chunked response body.
*   **Custom Error Handling**: Demonstrates 400 (Bad Request) and 500 (Internal Server Error) responses.

## Getting Started

### Prerequisites

*   Go (version 1.x recommended)

### Running the Server

1.  Clone the repository:
    ```bash
    git clone <your-repository-url>
    cd httpfromtcp
    ```
2.  Run the server:
    ```bash
    go run cmd/httpserver/main.go
    ```
    The server will start on port `42069` by default.

### Example Endpoints

Once the server is running, you can try accessing the following endpoints using a tool like `curl` or your web browser:

*   `http://localhost:42069/` - Returns a generic 200 OK HTML page.
*   `http://localhost:42069/yourproblem` - Returns a 400 Bad Request HTML page.
*   `http://localhost:42069/myproblem` - Returns a 500 Internal Server Error HTML page.
*   `http://localhost:42069/video` - Serves the `assets/vim.mp4` video file. (Make sure this file exists in an `assets` directory at the project root).
*   `http://localhost:42069/httpbin/get` - Proxies the request to `https://httpbin.org/get` and returns the response using chunked transfer encoding and trailers.
*   `http://localhost:42069/httpbin/headers` - Proxies to `https://httpbin.org/headers`.

## Project Structure

*   `cmd/httpserver/main.go`: Entry point of the application, sets up the server and request handlers.
*   `internal/server/`: Contains the core server logic for listening and handling connections.
*   `internal/request/`: Logic for parsing incoming HTTP requests.
*   `internal/response/`: Logic for constructing and writing HTTP responses, including status lines, headers, and body.
*   `internal/headers/`: Helper package for managing HTTP headers.
*   `assets/`: (Not version controlled by default - see `.gitignore`) Intended for static assets like the example video.

## Notes

*   The `assets/vim.mp4` file is expected to be in an `assets` folder at the root of the project for the `/video` endpoint to work. This folder is currently in `.gitignore`.
    *  To test this endpoint, you can create an `assets` directory and place a sample video at your choice. Make sure to adjust the path in the code if necessary.
*   This is an educational project and may not implement all aspects of the HTTP/1.1 specification or include robust error handling for all edge cases.