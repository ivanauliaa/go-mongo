package main

import (
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type student struct {
	Name  string `bson:"name"`
	Grade int    `bson:"grade"`
}

func Insert() {
	db, err := Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	students := []student{
		{
			Name:  "Wick",
			Grade: 2,
		}, {
			Name:  "Ethan",
			Grade: 2,
		},
	}

	for _, value := range students {
		_, err := db.Collection("student").InsertOne(GlobalCtx, value)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	fmt.Println("Insert success")
}

func Find() {
	db, err := Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	csr, err := db.Collection("student").Find(GlobalCtx, bson.M{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(GlobalCtx)

	result := []student{}
	for csr.Next(GlobalCtx) {
		row := student{}
		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	fmt.Println(result)
}

func Update() {
	db, err := Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	selector := bson.M{
		"name": "Wick",
	}

	changes := bson.M{
		"$set": student{
			Name:  "John Wick",
			Grade: 2,
		},
	}

	_, err = db.Collection("student").UpdateOne(GlobalCtx, selector, changes)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Update success")
}

func Delete() {
	db, err := Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	selector := bson.M{
		"name": "Wick",
	}

	_, err = db.Collection("student").DeleteOne(GlobalCtx, selector)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Delete success")
}

// Aggregate example
func Aggregate() {
	db, err := Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	pipeline := []bson.M{}
	err = bson.UnmarshalExtJSON([]byte(strings.TrimSpace(`
		[
			{
				"$group": {
					"_id": null,
					"Total": {
						"$sum": 1
					}
				}
			},
			{
				"$project": {
					"Total": 1,
					"_id": 0
				}
			}
		]
	`)), true, &pipeline)
	if err != nil {
		log.Fatal(err.Error())
	}

	csr, err := db.Collection("student").Aggregate(GlobalCtx, pipeline)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csr.Close(GlobalCtx)

	result := []bson.M{}

	for csr.Next(GlobalCtx) {
		row := bson.M{}

		err := csr.Decode(&row)
		if err != nil {
			log.Fatal(err.Error())
		}

		result = append(result, row)
	}

	fmt.Println(result)
}
