package main

import (
	"fmt"
	"log"

	"30.8/pkg/storage"
)

var db storage.Interface

func main() {

	var err error
	connstr := "host=localhost port=5432 user=postgres password=72evideb dbname=30.8 sslmode=disable"

	db, err := storage.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	//Новая таска
	/* task := storage.Task{
		Title:   "one",
		Content: "cont",
	} */
	/* id, err := db.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id) */

	//Удаление таски
	/* err = db.DeleteTask(7)
	if err != nil {
		log.Fatal(err)
	} */

	//Обновление таски
	/* task := storage.Task{
		AssignedID: 3,
		Title:      "one",
		Content:    "cont",
	}
	err = db.UpdateTask(8, task)
	if err != nil {
		log.Fatal(err)
	} */

	//Таска по id автора
	/* name, err := db.TaskByAuthor(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name) */

	tasks, err := db.Tasks(1, 2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasks)
}
