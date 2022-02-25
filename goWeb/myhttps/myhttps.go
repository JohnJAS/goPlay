package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

type PageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {

	// Make and parse the HTML template
	t, err := template.ParseFiles("todo.gohtml")
	if err != nil {
		panic(err)
	}

	// Initialze a struct storing page data and todo data
	data := PageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Render the data and output using standard output
		t.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}
