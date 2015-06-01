package main
import (
    "code.google.com/p/gorest"
    "net/http"
    "bufio"
    "bytes"
    "errors"
	"strings"
    "fmt"
    "io"
    "os"
    "encoding/json"
    //"container/list"
)


// Reprezentacja struktury odwróconego indeksu
// mapowanie string na int
var index map[string][]int
//tablica struktur doc czyli dokumentow
var indexed []doc

type SearchFiles struct {
    ID   int
    Filename string
    Folder string
   	Positions []positions
}

type positions struct {
	Line int
	Position int
}

//Struktura reprezentujaca dokument 
type doc struct {
    file  string
    title string
}

func main() {
    gorest.RegisterService(new(ElasticSearchService)) 
    http.Handle("/",gorest.Handle())    
    http.ListenAndServe(":4730",nil)
}

type ElasticSearchService struct {
    gorest.RestService 					`root:"/" consumes:"application/json" produces:"application/json"`
    getFilesList  gorest.EndPoint 		`method:"GET" path:"/files/" output:"string"`
    elasticSearch    gorest.EndPoint 	`method:"GET" path:"/search/{word:string}" output:"[]SearchFiles"`
	uploadFile	gorest.EndPoint 		`method:"GET" path:"/push/" output:"string"`
}

func(serv ElasticSearchService) GetFilesList() string{
    return "GetFilesList"
}
func(serv ElasticSearchService) ElasticSearch(word string) []SearchFiles{

	testArray:=[]SearchFiles{SearchFiles{}}

	//tworzymy mape za pomoca make
	index = make(map[string][]int)
	
	if err := indexDir("docs"); err != nil {
		testArray[0].ID=100
		testArray[0].Filename="errorsy jakies"
		testArray[0].Folder=" "
		
		b, _:= json.Marshal(testArray)
		n:= []SearchFiles{SearchFiles{}}
		json.Unmarshal(b,&n)
		return n
	}
	
	switch dl := index[word]; len(dl) {
		case 0:
			testArray[0].ID=0
			testArray[0].Filename="No match"
			testArray[0].Folder="docs"
		
			b, _:= json.Marshal(testArray)
			n:= []SearchFiles{SearchFiles{}}
			json.Unmarshal(b,&n)
			return n
		case 1:
			testArray[0].ID=0
			testArray[0].Filename=indexed[dl[0]].file
			testArray[0].Folder="docs"
			
			b, _:= json.Marshal(testArray)
			n:= []SearchFiles{SearchFiles{}}
			json.Unmarshal(b,&n)
			return n
		default:
			testArray2:=make([]SearchFiles,len(dl))

			for i:=0;i<len(dl);i++ {
				testArray2[i].ID=i
				testArray2[i].Filename=indexed[dl[i]].file
				testArray2[i].Folder="docs"

			}

			b, _:= json.Marshal(testArray2)
			n:= []SearchFiles{SearchFiles{}}
			json.Unmarshal(b,&n)
			return n
	}
	
}
func(serv ElasticSearchService) UploadFile() string{
    return "Upload file"
}

//Funkcja indeksujaca podany katalog
func indexDir(dir string) error {
	//otwieramy katalog do czytania
    df, err := os.Open(dir)

    if err != nil {
        return err
    }

	//podajac funkcji -1 zwracamy wszystkie FileInfo w jednym slice	
    fis, err := df.Readdir(-1)
    if err != nil {
        return err
    }
	//jesli slice == 0 to nie ma plikow w podanym katalogu
    if len(fis) == 0 {
        return errors.New(fmt.Sprintf("no files in %s", dir))
    }
	
	//ilosc zaindeksowanych plikow
    indexed := 0

	//iterujemy po slice fis
	//pomijamy iterator, potrzebna tylko wartosc
    for _, fi := range fis {
		//jesli nie jest katalogiem
        if !fi.IsDir() {
			//indeksujemy plik
            if indexFile(dir + "/" + fi.Name()) {
                //zwiekszamy indeks
				indexed++
				fmt.Println(indexed)
            }
        }
    }
    return nil
}

//Tworzenie odwroconego indeksu z pliku
func indexFile(filename string) bool {
    
	//otwieramy plik do czytania
	file, err := os.Open(filename)
    
	//jesli err nie jest nullem to zwracamy false i wypisujemy error
	if err != nil {
        fmt.Println(err)
        return false 
    }
 
	//jesli wszystko ok to rejestrujemy nowy plik
	//do zmiennej x przypisujemy dlugosc slice'a (wycinka)
    x := len(indexed)
	//do slice'a dodajemy slice i nowa strukture reprezentujaca dokument
    indexed = append(indexed, doc{filename, filename})
	//
    pdoc := &indexed[x]
 
    // robimy readera do skanowania pliku
    reader := bufio.NewReader(file)
	//ilosc linii
    lines := 0
	
	//for caly czas true - trwa do konca programu
    for {
		//Readline zwraca cala linie bez znaku konca linii
		//jesli linia jest za dluga to isPrefix jest ustawiany na true
		//zwracany jest wtedy poczatek linii i nastepna czesc linii bedzie przeczytana w nastepnym wywolaniu
		//buffer jest VALID do nastepnego wywolania ReadLine
		//bajty, isPrefix, error
        b, isPrefix, err := reader.ReadLine()
        switch {
		//jesli error == eof to true
        case err == io.EOF:
            return true
		//jesli error jest rozny od nulla to wypisujemy i zwracamy true
        case err != nil:
            fmt.Println(err)
            return true
		//jesli isPrefix jest ustawiony na true to zwracamy zbyt dluga linia
        case isPrefix:
            fmt.Printf("%s: Zbyt dluga linia w pliku.\n", filename)
            return true
		//jesli liczba linii jest mniejsza od 20
		//i bytes ma Title to zapisujemy document header
        case lines < 20 && bytes.HasPrefix(b, []byte("Title:")):
            pdoc.title = string(b[7:])
        }
		lines = lines + 1;
		//fmt.Println(lines)
		//indeksowana linia w tekscie w bajtach
		//trzeba poprawic word splitter'a
    wordLoop:
		//pomijamy  indeks w petli, potrzebne tylko wartosc, czyli bword
        for _, bword := range bytes.Fields(b) {
			//rozdzielamy bword za pomoca TRIM
            bword := bytes.Trim(bword, ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@ ")
            if len(bword) > 0 {
				//zapisujemy bword jako string do word
                word := string(bword)
				word = strings.ToLower(word)
				//nie mam pojecia co to robi????
                dl := index[word]
				//fmt.Println(dl)
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