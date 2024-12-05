package main

import (
	"log"
	"os"
)

func (api Api) show() {

	log.Println(api.FileName)
	log.Println(api.name)
	log.Println(api.Product)
	log.Println("Request:")
	for _, p := range api.Requests {
		log.Println(p.Name, "  ", p.Type)
	}
	log.Println()
	log.Println("Response:")
	for _, p := range api.Responses {
		log.Println(p.Name, "  ", p.Type)
	}
	log.Println("---------------------------------------")
	log.Println()

}
func show(apis []*Api) {
	log.Println(len(apis))
	for _, api := range apis {
		log.Println(api.FileName)
		log.Println("API Name:", api.Name)
		log.Printf("Product: %s, Version: %s, Name: %s\n", api.Product, api.Version, api.name)
		log.Println("Requests:")
		for _, req := range api.Requests {
			log.Printf("  Name: %s, Type: %s\n", req.Name, req.Type)
		}
		log.Println("Responses:")
		for _, res := range api.Responses {
			log.Printf("  Name: %s, Type: %s\n", res.Name, res.Type)
		}
		log.Println()
	}
}
func logInit() {
	file, err := os.Create("output/log.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

}
