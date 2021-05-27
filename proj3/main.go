package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

func TFWorker(files []fs.FileInfo, start int, end int, vectors []*map[string]float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i < end; i++ {
		file := files[i]
		dat, _ := ioutil.ReadFile("books/" + file.Name())

		s := strings.ReplaceAll(string(dat), "\n", " ")
		ss := strings.Split(s, " ")

		d := make(map[string]float64)
		for _, v := range ss {
			d[v] += 1.0
		}

		vectors[i] = &d
	}
}

func IDFWorker(words []string, vectors []*map[string]float64, df []float64, idf []float64, start int, end int, n_books int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i < end; i++ {
		word := words[i]
		for _, p := range vectors {
			if (*p)[word] > 0 {
				df[i] += 1.0
			}
		}
		idf[i] = math.Log(float64(n_books) / df[i])
	}
}

func TFIDFWorker(words []string, vectors []*map[string]float64, df []float64, idf []float64, start int, end int, n_books int, tfidf [][]float64, wg *sync.WaitGroup) {
	defer wg.Done()
	//For every document
	for b := start; b < end; b++ {
		//Make an array of length 10000
		this_book := make([]float64, len(words))
		//Calculate the TF-IDF for each word
		for i := 0; i < len(words); i++ {
			this_book[i] = (*vectors[b])[words[i]] * idf[i]
		}
		//Append the result
		tfidf[b] = this_book
		if b == 100 {
			fmt.Println(this_book)
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("ERROR: Wrong number of command line arguments")
		fmt.Println("Usage ./bagOfWords <n_threads> <n_books>")
		fmt.Println("n_threads of 1 will run a serial version")
		return
	}
	n_threads, _ := strconv.Atoi(os.Args[1])
	n_books, _ := strconv.Atoi(os.Args[2])

	//Read the 10000 most common words
	file, err := os.Open("google-10000-english.txt")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	//A vector to hold pointers to our maps
	vectors := make([]*map[string]float64, n_books)
	df := make([]float64, len(words))
	idf := make([]float64, len(words))
	// var tfidf [][]float64
	tfidf := make([][]float64, n_books)

	//Get our directory of books
	files, err := ioutil.ReadDir("books")
	if err != nil {
		panic("ERROR: Cannot find books directory.")
	}
	if n_threads == 1 { /************** SERIAL VERSION **************/
		//Calculate term frequencies for each book
		for i := 0; i < n_books; i++ {
			file := files[i]
			dat, _ := ioutil.ReadFile("books/" + file.Name())

			s := strings.ReplaceAll(string(dat), "\n", " ")
			ss := strings.Split(s, " ")

			d := make(map[string]float64)
			for _, v := range ss {
				d[v] += 1.0
			}

			vectors[i] = &d
		}

		//Calculate Inverse Document Frequency for each word
		for i := 0; i < len(words); i++ {
			word := words[i]
			for _, p := range vectors {
				if (*p)[word] > 0 {
					df[i] += 1.0
				}
			}
			idf[i] = math.Log(float64(n_books) / df[i])
		}

		//Calculate TFIDF for every word in each document

		//For every document
		for b := 0; b < n_books; b++ {
			//Make an array of length 10000
			this_book := make([]float64, len(words))
			//Calculate the TF-IDF for each word
			for i := 0; i < len(words); i++ {
				this_book[i] = (*vectors[b])[words[i]] * idf[i]
			}
			//Append the result
			tfidf[b] = this_book
			if b == 100 {
				fmt.Println(this_book)
			}
		}

	} else { /************** PARALLEL VERSION **************/
		var wg sync.WaitGroup

		work := n_books / n_threads
		for i := 0; i < n_threads; i++ {
			wg.Add(1)
			go TFWorker(files, i*work, (i+1)*work, vectors, &wg)
		}
		wg.Wait()

		work2 := len(words) / n_threads
		for i := 0; i < n_threads; i++ {
			wg.Add(1)
			go IDFWorker(words, vectors, df, idf, i*work2, (i+1)*work2, n_books, &wg)
		}
		wg.Wait()

		for i := 0; i < n_threads; i++ {
			wg.Add(1)
			go TFIDFWorker(words, vectors, df, idf, i*work, (i+1)*work, n_books, tfidf, &wg)
		}
		wg.Wait()

	}
}
