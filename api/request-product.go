package main

import (
	"encoding/json"
	"log"
	"net/http"

	"./product"
	"./user"
	"github.com/gorilla/mux"
)

func getProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	barcode := params["barcode"]
	product := product.GetProductInfo(barcode, conn)
	log.Printf("Barcode %s with an error of? %s", barcode, product.Error)
	if product.Error != "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	username := getUsername(r)
	user.PointsForScan(product.ID, product.Version, username, conn)
	output, _ := json.Marshal(product)
	w.Write(output)

}

func changeProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p product.Product
	err := decoder.Decode(&p)
	failOnError(err, "Failed to decode product")
	product.AlterProduct(p, getUsername(r), conn)
}
