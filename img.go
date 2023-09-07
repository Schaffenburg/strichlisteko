package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"

	"image"
	_ "image/jpeg"
	_ "image/png"
)

func handlePutImage(w http.ResponseWriter, r *http.Request) {
	var uuid = uuid.New()

	// read image from req body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body for uuid: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, "{\"error\":\"failed to read body.\"}")
		return
	}

	_, mime, err := image.Decode(bytes.NewReader(b))
	if err != nil && mime != "" {
		log.Printf("failed to decode image: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, "{\"error\":\"failed to decode image.\"}")
		return
	}

	err = setImage(uuid, mime, b)
	if err != nil {
		log.Printf("failed to put image into db: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, "{\"error\":\"failed to write image to db.\"}")
		return
	}

	fmt.Fprintf(w, "{\"info\":\"upload %s; len %d; mime %s.\",\"uuid\":\"%s\"}",
		uuid.String(), len(b), mime,
		uuid.String(),
	)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", "/img/"+uuid.String())
	w.WriteHeader(201)

	fmt.Fprintf(w, "Wrote %d bytes to db. uuid is %s. mime is %s\n", len(b), uuid, mime)
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr, ok := vars["image"]
	if !ok {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(404)
		w.Write(image404)
		return
	}

	uuid := new(uuid.UUID)
	err := uuid.UnmarshalText([]byte(idstr))
	if err != nil {
		log.Printf("failed to parse imguuid '%s': %s", idstr, err)
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(404)
		w.Write(image404)
		return
	}

	img, mime, err := getImage(*uuid)
	if err != nil {
		log.Printf("Failed to get image %s: %s", uuid, err)
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(404)
		w.Write(image404)
		return
	}

	w.Header().Set("Content-Type", "image/"+mime)
	w.WriteHeader(200)
	w.Write(img)
}

func handlePutImageAPI(w http.ResponseWriter, r *http.Request) {
	var uuid = uuid.New()

	// read image from req body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body for uuid: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, "{\"error\":\"failed to read body.\"}")
		return
	}

	_, mime, err := image.Decode(bytes.NewReader(b))
	if err != nil && mime != "" {
		log.Printf("failed to decode image: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, "{\"error\":\"failed to decode image.\"}")
		return
	}

	err = setImage(uuid, mime, b)
	if err != nil {
		log.Printf("failed to put image into db: %s", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		fmt.Fprint(w, "{\"error\":\"failed to write image to db.\"}")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	fmt.Fprintf(w, "{\"bytes\": %d, \"uuid\":\"%s\", \"mime\": \"%s\", \"info\": \"success\"}", len(b), uuid, mime)
}
