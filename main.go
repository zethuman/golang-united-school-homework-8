package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func unmarshalAll(data []byte) []Item {
	var items []Item
	json.Unmarshal(data, &items)
	return items
}

func marshalAll(items []Item) []byte {
	datas, _ := json.Marshal(items)
	return datas
}

func marshalOne(item Item) []byte {
	data, _ := json.Marshal(item)
	return data
}

func unmarshalOne(data string) Item {
	var item Item
	json.Unmarshal([]byte(data), &item)
	return item
}

func getItems(fileName string) []Item {
	file, _ := os.Open(fileName)
	list, _ := ioutil.ReadAll(file)
	unmarshalledAll := unmarshalAll(list)
	return unmarshalledAll
}

func write(fileName string, items []Item) {
	marshalled := marshalAll(items)
	ioutil.WriteFile(fileName, []byte(marshalled), 0666)
}

func list(fileName string, w io.Writer) {
	items := getItems(fileName)
	w.Write(marshalAll(items))
}

func add(item string, fileName string, w io.Writer) {
	items := getItems(fileName)
	unmarshalledOne := unmarshalOne(item)
	for _, v := range items {
		if v.Id == unmarshalledOne.Id {
			w.Write([]byte(fmt.Errorf("Item with id %s already exists", v.Id).Error()))
		}
	}
	items = append(items, unmarshalledOne)
	write(fileName, items)
}

func remove(id string, fileName string, w io.Writer) {
	items := getItems(fileName)
	var found bool
	for i, v := range items {
		if v.Id == id {
			found = true
			items = append(items[:i], items[i+1:]...)
		}
	}

	if !found {
		w.Write([]byte(fmt.Errorf("Item with id %s not found", id).Error()))
	}

	write(fileName, items)
}

func findById(id string, fileName string, w io.Writer) {
	items := getItems(fileName)
	var found bool
	for _, v := range items {
		if v.Id == id {
			found = true
			w.Write(marshalOne(v))
		}
	}

	if !found {
		w.Write([]byte(""))
	}
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}

	if args["fileName"] == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	if args["operation"] != "add" && args["operation"] != "remove" && args["operation"] != "findById" && args["operation"] != "list" {
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}

	switch args["operation"] {
	case "list":
		list(args["fileName"], writer)
	case "add":
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		add(args["item"], args["fileName"], writer)
	case "remove":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		remove(args["id"], args["fileName"], writer)
	case "findById":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		findById(args["id"], args["fileName"], writer)
	}

	return nil
}

func parseArgs() Arguments {
	var operation, fileName, item string

	flag.StringVar(&operation, "operation", "", "./main -operation \"add\"")
	flag.StringVar(&fileName, "fileName", "", "./main -fileName \"add\"")
	flag.StringVar(&item, "item", "", "./main -item \"add\"")

	flag.Parse()

	return Arguments{
		"operation": operation,
		"fileName":  fileName,
		"item":      item,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
