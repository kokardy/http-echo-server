# http-echo-server

A simple HTTP echo server written in Go, inspired by `httpbin.org`.

The server displays the request body as streaming chunks, making it memory efficient.

## Installation

First, ensure you have Go installed and set up on your machine. Then, you can install the http-echo-server using the following command:

```bash
$ go install github.com/kokardy/http-echo-server/cmd/http-echo-server@latest
```

## Usage

To start the server, run:

```bash
$ http-echo-server --port 10080
```

## Examples

The examples use the `httpie` command, but feel free to use `curl` if you prefer.

### Using HTTPie

```bash
$ http -v --form 'localhost:10080/fff?q=100' a=b c=d
```

### Using cURL

```bash
$ curl -v -X POST 'localhost:10080/fff?q=100' -d 'a=b&c=d'
```

## Response

The server will respond with the request details in JSON format:

```json
{
    "body_chunks": [
        {
            "chunk_num": 0,
            "data": "a=b&c=d",
            "size": 7
        }
    ],
    "headers": {
        "Accept": "*/*",
        "Accept-Encoding": "gzip, deflate",
        "Connection": "keep-alive",
        "Content-Length": "7",
        "Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
        "Host": "localhost:10080",
        "User-Agent": "HTTPie/2.6.0"
    },
    "metadata": {
        "client_ip": "127.0.0.1",
        "method": "POST",
        "timestamp": "2025-03-02T15:40:31+09:00"
    },
    "path": "/fff",
    "query_params": {
        "q": [
            "100"
        ]
    }
}
```

## License

This project is licensed under the MIT License.
