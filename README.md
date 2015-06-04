# ElasticSearchGo
Elastic search with Go Lang - parallel REST API.

# Ubuntu
Install package:
```
sudo apt-get install golang golang-go.tools gccgo-go git mercurial
```

Set variables:
```
export GOROOT=/usr/lib/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
Install github package:
```
go get -u github.com/gorrila/mux
```
Open browser and test:
```
localhost:4730/files
localhost:4730/search/is
```
cURL for upload sample file:
```
curl -F file=@SCIEZKA/DO/PLIKU http://localhost:4730/push
```