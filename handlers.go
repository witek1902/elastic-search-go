package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

func Files(w http.ResponseWriter, r *http.Request) {
	
	if err := initDocumets("docs"); err != false {

		tempFileArray := make([]File,len(Documents))			
		
		for i:=0;i<len(Documents);i++ {
			tempFileArray[i].Id=i
			tempFileArray[i].Filename=Documents[i].File 
			tempFileArray[i].Folder="docs"
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		
		if err := json.NewEncoder(w).Encode(tempFileArray); err != nil {
			panic(err)
		}
		return
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	
	word := vars["word"]

	fileArray := []SearchFile{SearchFile{}}
	index = make(map[string][]int)
	
	if err := indexDir("docs"); err != nil {
		fileArray[0].Id=0
		fileArray[0].Filename="Dir not exists."
		fileArray[0].Folder="docs"
		
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		
		if err := json.NewEncoder(w).Encode(fileArray); err != nil {
			panic(err)
		}
		return
	}

	//tablica intow pasujaca dla danego slowa
	switch documentList := index[word]; len(documentList) {
		case 0:
			fileArray[0].Id=0
			fileArray[0].Filename="No match files"
			fileArray[0].Folder="docs"
			
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
		
			if err := json.NewEncoder(w).Encode(fileArray); err != nil {
				panic(err)
			}
			return
		case 1:
			fileArray[0].Id=0
			fileArray[0].Filename=Documents[documentList[0]].File
			fileArray[0].Folder="docs"
			
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(fileArray); err != nil {
				panic(err)
			}
			return
		default:
			tempFileArray := make([]SearchFile,len(documentList))			
			
			for i:=0;i<len(documentList);i++ {
				tempFileArray[i].Id=i
				tempFileArray[i].Filename=Documents[documentList[i]].File
				tempFileArray[i].Folder="docs"
			}
			
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(tempFileArray); err != nil {
				panic(err)
			}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}


/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

*/
func Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
	var fileUpload fileUpload
	body,err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &fileUpload); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	
//	t := RepoCreateTodo(todo)
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusCreated)
//	if err := json.NewEncoder(w).Encode(t); err != nil {
//		panic(err)
//	}
}

