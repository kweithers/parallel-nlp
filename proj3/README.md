# Project 3: Parallel Natural Language Processing

## Overview

A ![bag-of-words](https://en.wikipedia.org/wiki/Bag-of-words_model) model can be used to create a feature vector that describes a document. It simply counts how many times each word occurs in a document.

For example, the document "John likes to watch movies. Mary likes movies too." has the vector ```{"John":1,"likes":2,"to":1,"watch":1,"movies":2,"Mary":1,"too":1}```.

For each document, we want to get a 10000 integer vector representing the ![most common 10,000 words](https://raw.githubusercontent.com/first20hours/google-10000-english/master/google-10000-english.txt) so that we can compare different documents and find similarities between them.

## Data 

I pulled 1000 books from ![Project Guttenberg](https://www.gutenberg.org/) for the dataset. I made a simple bash script ```books/get_books.sh``` that calls wget a bunch of times to download the books. Note there are two commonly used url formats, so I just try both for each book_id.

While getting a bag-of-words vector as quickly as possible might not be super important for books that are decades old, one could easily imagine applying this very same technique to tweets that are coming in real time, scraping news stories from major news orgs, or many other uses cases when speed is paramount.

## Parallel bag-of-words vector generation.

This problem is embarassingly parallel. Each book can be assigned to a goroutine that calculates its bag-of-words vector completely independently of every other book, so this makes a perfect use case for parallel programming and we should see some nice speedup!

## Term frequency-inverse document frequency 

![Term frequency-inverse document frequency (tf-idf)](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) takes this one step further. It is a slightly more sophisticated statistic that can be used to represent how important certain words are to a document. Instead of simply counting how many times a word appears, we include another term called *inverse document frequency*. It accounts for what percentage of *documents* a word appears in as well. This way common words like 'the' and 'an' don't get a big weight since they appear in every document, but rarer words get a bigger weight since they are more likely to be important.

To do this, we need to wait until all of the workers are finished counting their individual documents, we do this with a wait group.

Once they are done, we spawn more workers that calculate, for each word, how many documents that word appears in. Now, we can calculate the inverse document frequency for each word.

Finally, we can spawn workers that calculate the tfidf for each word in each document by multiplying the term frequency and the inverse document frequency.

Now for each document we have a 10000 feature vector. We can use this for tons of things, like finding the ![cosine similarity](https://en.wikipedia.org/wiki/Cosine_similarity) between two documents feature vectors as a measure of their similariy, or using it as an input for some machine learning model.

## How to run the program

0. cd into the proj3 directory

1. cd into the books directory, then call ``` bash ./get_books.sh```

2. cd back into the proj3 directory, then call ```go run main.go <n_threads> <n_books> <save_results>```
  * n_threads: the number of worker goroutines to create. if this is equal to 1, run the serial version
  * n_books: how many of our 1000 books to run the process on (note: ensure that n_books / n_threads is an integer)
  * save_reults: 0 if you don't want to save results, 1 if you want to save results into the tfidf folder

3. if you want to produce your own speedup graph, call ```python timings.py```

## Performance Discussion

![](Speedup.png)

I ran these benchmarks on my 2019 MacBook Pro with the following specs

* Processor: 2.6 GHz 6-Core Intel Core i7
* RAM: 16 GB 2667 MHz DDR4

The number of documents doesn't seem to have a huge impact on the speedup itself. However, the speedup flattens out after 6, which is equal to the number of cores on my machine, which makes sense. Even though hyperthreading is enabled, it seems to not have an effect here.

Parallel programming was able to provide some nice speeedup for this use case!