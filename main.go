package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/selfidrone/detection/faces"
)

var log hclog.Logger

func main() {
	log = hclog.Default()
	log.Info("Starting server")

	http.HandleFunc("/detect", handle)
	http.HandleFunc("/health", health)
	http.ListenAndServe(":9999", nil)
}

// Request is base64 encoded image

// Response for the function
type Response struct {
	Faces       []image.Rectangle
	Bounds      image.Rectangle
	ImageBase64 string
}

func health(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

// Handle a serverless request
func handle(rw http.ResponseWriter, r *http.Request) {
	defer func(t time.Time) {
		log.Info("Finished processing request", "time", time.Now().Sub(t))
	}(time.Now())

	var body []byte
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("No body received")

		http.Error(rw, "", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Error("Unable to decode base64")

		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	typ := http.DetectContentType(data)
	if typ != "image/jpeg" && typ != "image/png" {
		log.Error("Image is not jpeg or png")

		http.Error(
			rw,
			`Only jpeg or png images, either raw uncompressed bytes or base64 encoded 
			are acceptable inputs, you uploaded: `+typ,
			http.StatusBadRequest,
		)
		return
	}

	tmpfile, err := ioutil.TempFile("/tmp", "image")
	if err != nil {
		log.Error("Unable to write file")

		http.Error(rw, "", http.StatusBadRequest)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	io.Copy(tmpfile, bytes.NewBuffer(data))

	faceProcessor := faces.NewFaceProcessor()
	faces, bounds := faceProcessor.DetectFaces(tmpfile.Name())

	resp := Response{
		Faces:  faces,
		Bounds: bounds,
	}

	// do we need to create and output an image?
	output := r.URL.Query().Get("output")
	var image []byte
	if output == "image" || output == "json_image" {
		var err error
		image, err = faceProcessor.DrawFaces(tmpfile.Name(), faces)
		if err != nil {
			log.Error("Unable to create image output", "error", err)

			http.Error(rw, fmt.Sprintf("Error creating image output: %s", err), http.StatusInternalServerError)
			return
		}

		resp.ImageBase64 = base64.StdEncoding.EncodeToString(image)
	}

	if output == "image" {
		log.Info("Output image")

		rw.Write(image)
		return
	}

	enc := json.NewEncoder(rw)
	enc.Encode(resp)
}
