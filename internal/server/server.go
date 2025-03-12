package server

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
}

// New creates a new server instance
func New() *Server {
	router := gin.Default()

	server := &Server{
		router: router,
	}

	server.registerRoutes()

	return server
}

// registerRoutes sets up the server routes
func (s *Server) registerRoutes() {
	s.router.Any("/*path", s.echoHandler)
}

// echoHandler handles all requests and streams back request details
func (s *Server) echoHandler(c *gin.Context) {
	// Set response headers
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusOK)

	// Get headers
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// Add Host header
	headers["Host"] = c.Request.Host

	// Add Transfer-Encoding header if present
	if len(c.Request.TransferEncoding) > 0 {
		headers["Transfer-Encoding"] = c.Request.TransferEncoding[0]
	}

	// Prepare response data
	responseData := gin.H{
		"method":    c.Request.Method,
		"client_ip": c.ClientIP(),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Determine the accepted encoding
	acceptEncoding := c.GetHeader("Accept-Encoding")
	var writer io.Writer = c.Writer
	var closer io.Closer
	var encoding string

	if strings.Contains(acceptEncoding, "br") {
		encoding = "br"
		brWriter := brotli.NewWriter(c.Writer)
		writer = brWriter
		closer = brWriter
	} else if strings.Contains(acceptEncoding, "gzip") {
		encoding = "gzip"
		gzipWriter := gzip.NewWriter(c.Writer)
		writer = gzipWriter
		closer = gzipWriter
	} else if strings.Contains(acceptEncoding, "deflate") {
		encoding = "deflate"
		zlibWriter := zlib.NewWriter(c.Writer)
		writer = zlibWriter
		closer = zlibWriter
	}

	if encoding != "" {
		defer closer.Close()
		c.Header("Content-Encoding", encoding)
		fmt.Printf("Response Content-Encoding: %s\n", encoding)
	}

	// Stream the response
	c.Stream(func(w io.Writer) bool {
		// Marshal the metadata part
		metaJSON, err := json.Marshal(responseData)
		if err != nil {
			return false
		}
		headersJSON, err := json.Marshal(headers)
		if err != nil {
			return false
		}
		queryParamsJSON, err := json.Marshal(c.Request.URL.Query())
		if err != nil {
			return false
		}
		fmt.Fprintf(writer, `{"metadata": %s, "path": "%s", "headers": %s, "query_params": %s, "body_chunks": [`, metaJSON, c.Request.URL.Path, headersJSON, queryParamsJSON)

		// If there's a body, stream it in chunks
		if c.Request.Body != nil {
			var reader io.ReadCloser
			contentEncoding := c.GetHeader("Content-Encoding")
			switch contentEncoding {
			case "gzip":
				reader, err = gzip.NewReader(c.Request.Body)
			case "deflate":
				reader, err = zlib.NewReader(c.Request.Body)
			case "br":
				reader = io.NopCloser(brotli.NewReader(c.Request.Body))
			default:
				reader = c.Request.Body
			}
			if err != nil {
				return false
			}
			defer reader.Close()

			// Log the content encoding
			fmt.Printf("Request Content-Encoding: %s\n", contentEncoding)

			buffer := make([]byte, 4096) // 4KB chunks
			bodyPartCount := 0
			firstChunk := true

			for {
				n, err := reader.Read(buffer)
				if n > 0 {
					if !firstChunk {
						fmt.Fprint(writer, ",")
					}
					firstChunk = false

					bodyPart := gin.H{
						"chunk_num": bodyPartCount,
						"data":      string(buffer[:n]),
						"size":      n,
					}

					bodyJSON, err := json.Marshal(bodyPart)
					if err != nil {
						break
					}

					fmt.Fprintf(writer, "%s", bodyJSON)
					bodyPartCount++
				}

				if err == io.EOF {
					// End of body
					break
				}

				if err != nil {
					// Error reading body
					errMsg := gin.H{
						"error": err.Error(),
					}
					errJSON, _ := json.Marshal(errMsg)
					fmt.Fprintf(writer, `,{"error": %s}`, errJSON)
					break
				}
			}
		}

		// Close the body_chunks array and the JSON object
		fmt.Fprint(writer, "]}")
		return false
	})
}

// Run starts the server on the given address
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
