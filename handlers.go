package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"github.com/gorilla/mux"
)

func Files(w http.ResponseWriter, r *http.Request) {
	
		tempFileArray1 := make([]File,len(Documents))			
		
		for i:=0;i<len(Documents);i++ {
			tempFileArray1[i].Id=i
			tempFileArray1[i].Filename=Documents[i].File 
			tempFileArray1[i].Folder="docs"
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		
		if err := json.NewEncoder(w).Encode(tempFileArray1); err != nil {
			panic(err)
		}
		return
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
curl -F file=@SCIEZKA/DO/PLIKU 
http://localhost:4730/push

*/
func Upload(w http.ResponseWriter, r *http.Request) {
 
	file, header, err := r.FormFile("file")
    if err != nil {
		fmt.Println("Can't find file")
		return 
    }

	temp, _ := ioutil.TempFile("./docs/", header.Filename+"-")
    defer temp.Close()

    _, err = io.Copy(temp, file)
    if err != nil {
            fmt.Fprintln(w, err)
    }

    fmt.Fprintf(w, "File uploaded successfully : ")
    fmt.Fprintf(w, header.Filename)
	 
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

