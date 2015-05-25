package main
import (
    "code.google.com/p/gorest"
    "net/http"
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "io"
    "os"
)

// Reprezentacja struktury odwróconego indeksu
var index map[string][]int
var indexed []doc

//Struktura reprezentujaca dokument 
type doc struct {
    file  string
    title string
}

func main() {
    gorest.RegisterService(new(ElasticSearchService)) 
    http.Handle("/",gorest.Handle())    
    http.ListenAndServe(":8787",nil)
}

type ElasticSearchService struct {
    gorest.RestService `root:"/"`
    getFilesList  gorest.EndPoint `method:"GET" path:"/files/" output:"string"`
    elasticSearch    gorest.EndPoint `method:"GET" path:"/search/{word:string}" output:"string"`
	uploadFile	gorest.EndPoint `method:"GET" path:"/push/" output:"string"`
}
func(serv ElasticSearchService) GetFilesList() string{
    return "GetFilesList"
}
func(serv ElasticSearchService) ElasticSearch(word string) string{
    var buffer bytes.Buffer
	
	index = make(map[string][]int)
	
	if err := indexDir("docs"); err != nil {
		return "Error"
	}
	
	switch dl := index[word]; len(dl) {
		case 0:
			return "No match"
		case 1:
			return indexed[dl[0]].file
		default:
			buffer.WriteString("Search word: ")
			buffer.WriteString(word)
			buffer.WriteString("\n")
			for _, d := range dl {
				buffer.WriteString(indexed[d].file)
				buffer.WriteString("  \n")
			}
			return buffer.String()
	}
	
}
func(serv ElasticSearchService) UploadFile() string{
    return "Upload file"
}

//Funkcja wyszukujaca katalog
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
    indexed := 0
    for _, fi := range fis {
        if !fi.IsDir() {
            if indexFile(dir + "/" + fi.Name()) {
                indexed++
            }
        }
    }
    return nil
}

//Tworzenie odwroconego indeksu z pliku
func indexFile(fn string) bool {
    f, err := os.Open(fn)
    if err != nil {
        fmt.Println(err)
        return false // only false return
    }
 
    // register new file
    x := len(indexed)
    indexed = append(indexed, doc{fn, fn})
    pdoc := &indexed[x]
 
    // scan lines
    r := bufio.NewReader(f)
    lines := 0
    for {
        b, isPrefix, err := r.ReadLine()
        switch {
        case err == io.EOF:
            return true
        case err != nil:
            fmt.Println(err)
            return true
        case isPrefix:
            fmt.Printf("%s: unexpected long line\n", fn)
            return true
        case lines < 20 && bytes.HasPrefix(b, []byte("Title:")):
            // in a real program you would write code
            // to skip the Gutenberg document header
            // and not index it.
            pdoc.title = string(b[7:])
        }
        // index line of text in b
        // again, in a real program you would write a much
        // nicer word splitter.
    wordLoop:
        for _, bword := range bytes.Fields(b) {
            bword := bytes.Trim(bword, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@")
            if len(bword) > 0 {
                word := string(bword)
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
    return true
}   