package main

var index map[string][]int

type SearchFile struct {
	Id int
	Filename string
	Folder	string
	Position []Position
}

type File struct {
	Id int
	Filename string
	Folder	string
}

type Position struct {
	Line int
	Position int
}

type Document struct {
	File string
	Title string
}
var Documents []Document

type fileUpload struct {
	filename string
	data []byte
	folder string
}
