package video

import (
	"net/http"
	"os"
	"strconv"
	"strings"
)

//
func ServeFromLocalFile(w http.ResponseWriter, r *http.Request) {
	// Open video file (& defer close)
	filepath := "./assets/sample_h264.mp4"
	file, err := os.Open(filepath)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer file.Close()

	// Read file info to get size
	fi, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get video size
	fileSize := int(fi.Size())

	// Define buffer size
	const bufferSize = 1024 * 8

	// Check if request header contains Range field
	rangeHeader := r.Header.Get("range")
	// Range is not defined yet means it is first request (beginning of video)
	if rangeHeader == "" {
		contentLength := strconv.Itoa(fileSize)
		contentEnd := strconv.Itoa(fileSize - 1)

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes 0-"+contentEnd+"/"+contentLength)
		w.WriteHeader(200)

		buffer := make([]byte, bufferSize)

		for {
			n, err := file.Read(buffer)

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			data := buffer[:n]
			w.Write(data)
			w.(http.Flusher).Flush()
		}
	} else {
		// Range is defined
		rangeParam := strings.Split(rangeHeader, "=")[1]
		splitParams := strings.Split(rangeParam, "-")

		contentStartValue := 0
		contentStart := strconv.Itoa(contentStartValue)
		contentEndValue := fileSize - 1
		contentEnd := strconv.Itoa(contentEndValue)
		contentSize := strconv.Itoa(fileSize)

		if len(splitParams) > 0 {
			contentStartValue, err = strconv.Atoi(splitParams[0])

			if err != nil {
				contentStartValue = 0
			}

			contentStart = strconv.Itoa(contentStartValue)
		}

		if len(splitParams) > 1 {
			contentEndValue, err = strconv.Atoi(splitParams[1])

			if err != nil {
				contentEndValue = fileSize - 1
			}

			contentEnd = strconv.Itoa(contentEndValue)
		}

		contentLength := strconv.Itoa(contentEndValue - contentStartValue + 1)

		// Set response Header
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Range", "bytes "+contentStart+"-"+contentEnd+"/"+contentSize)
		w.WriteHeader(206)

		// Write response data
		buffer := make([]byte, bufferSize)

		_, err := file.Seek(int64(contentStartValue), 0)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeBytes := 0

		for {
			n, err := file.Read(buffer)

			writeBytes += n

			if n == 0 {
				break
			}

			if err != nil {
				break
			}

			if writeBytes >= contentEndValue {
				data := buffer[:bufferSize-writeBytes+contentEndValue+1]
				w.Write(data)
				w.(http.Flusher).Flush()
				break
			}

			data := buffer[:n]
			w.Write(data)
			w.(http.Flusher).Flush()
		}
	}
}
