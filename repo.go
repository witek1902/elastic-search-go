package main

import (
	"fmt"
	"os"
	"errors"
	"bufio"
	"io"
	"bytes"
	"strings"
	"log"
)


func initDocumets(dir string) bool {
	df, err := os.Open(dir)
    if err != nil {
        return false
    }

    fis, err := df.Readdir(-1)
    if err != nil {
        return false
    }
	
    for _, fi := range fis {
        if !fi.IsDir() {
            Documents = append(Documents, Document{fi.Name(),fi.Name()})
        }
    }
    return true
}

func indexDir(dir string) error {
	
    df, err := os.Open(dir)
    if err != nil {
        return err
    }

    fis, err := df.Readdir(-1)
    if err != nil {
        return err
    }
	
    if len(fis) == 0 {
        return errors.New(fmt.Sprintf("no files in %s", dir))
    }
	
    for _, fi := range fis {
        if !fi.IsDir() {
            indexFile(dir + "/" + fi.Name())
        }
    }
    return nil
}

func indexFile(filename string) bool {
    
	file, err := os.Open(filename)
	if err != nil {
        fmt.Println(err)
        return false 
    }
 
	//x = liczba dokumentow
    x := len(Documents)
	//do Documents dodajemy kolejny dokument
    Documents = append(Documents, Document{filename, filename})
 
	//wczytujemy kawalki
    reader := bufio.NewReader(file)
	nBytes, nChunks := int64(0), int64(0)
	buf := make([]byte, 0,4*1024)
    for {
		n, err := reader.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
            if err == nil {
                continue
            }
            if err == io.EOF {
                break
            }
            log.Fatal(err)
        }
		
		nChunks++
		nBytes += int64(len(buf))
        if err != nil && err != io.EOF {
            log.Fatal(err)
        }
	
	//wczytujemy kolejne slowa 
    wordLoop:
        for _, bword := range bytes.Fields(buf) {
            bword := bytes.Trim(bword, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@ ")
            if len(bword) > 0 {
                word := string(bword)
				word = strings.ToLower(word)
				
                dl := index[word]
				
                for _, d := range dl {
                    if d == x {
                        continue wordLoop
                    }
                }
                index[word] = append(dl, x)
            }
        }
    }
	log.Println("Bytes:", nBytes, "Chunks:", nChunks)
    return true
}   
