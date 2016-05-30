package utils

//package main

import (
	"bufio"
	//"fmt"
	//"io"
	"os"
	//"time"
)

type fileScanner struct {
	File    *os.File
	Scanner *bufio.Scanner
	Reader  *bufio.Reader
}

func NewFileScanner(fileName string) (*fileScanner, error) {
	tmpFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fs := fileScanner{File: tmpFile}
	return &fs, nil
}

func (f *fileScanner) Close() error {
	return f.File.Close()
}

/*func (f *fileScanner) GetReader() *bufio.Reader {
	if f.Reader == nil {
		f.Reader = bufio.NewReader(f.File)
	}
	return f.Reader
}*/

func (f *fileScanner) GetScanner() *bufio.Scanner {
	if f.Scanner == nil {
		f.Scanner = bufio.NewScanner(f.File)
		//f.Scanner.Split(bufio.ScanLines)
	}
	return f.Scanner
}

/*func main() {
	var path string = "/home/k/work/go/src/github.com/go-stomp/stomp/examples/client_test/test.csv"
	//if len(os.Args) > 1 {
	fmt.Println("11")
	fscanner, err := NewFileScanner(path)
	//err := fscanner.Open(path)
	if err == nil {
		defer fscanner.Close()
		fmt.Println("22")
		scanner := fscanner.GetScanner()
		fmt.Println("33")
		// Go through file line by line.
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			time.Sleep(1 * time.Second)
			// or do other stuff
		}

	} else {
		fmt.Println("errr")
	}
	//}
}*/
