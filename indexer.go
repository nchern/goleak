package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	//	"runtime"
	"runtime/pprof"
)

var (
	inputFile = flag.String("infile", "enwiki-latest-pages-articles.xml", "Input file path")
)

func pageProcessor(queue <-chan *Page, results chan<- int) {
	for page := range queue {
		page.Text += ""
		results <- 1
	}
}

func main() {
	//runtime.GOMAXPROCS(2)

	flag.Parse()

	xmlFile, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("Error opening file:", err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.Entity = xml.HTMLEntity
	total := 0
	maxTotal := 350000000
	var inElement string

	queue := make(chan *Page, 100000)
	results := make(chan int, 100000)

	var size int64

	go pageProcessor(queue, results)

	for {
		token, _ := decoder.Token()
		if token == nil || total > maxTotal {
			break
		}
		switch tokenType := token.(type) {
		case xml.StartElement:
			inElement = tokenType.Name.Local
			if inElement == "page" {
				var p Page
				decoder.DecodeElement(&p, &tokenType)

				queue <- &p

				total++
				size += int64(len(p.Text))
				if total%10000 == 0 {
					pprof.Lookup("heap").WriteTo(os.Stderr, 1)
					log.Printf("Processed %d, total size: %v", total, size)
				}
				//if total%10000 == 0 {
				//	runtime.GC()
				//}
			}
		default:
		}
	}
	fmt.Printf("Total read articles: %d \n", total)
	close(queue)
	for i := 0; i < total; i++ {
		<-results
	}
}
