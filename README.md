# Project 3: Parallel Natural Language Processing

## Overview

A ![bag-of-words](https://en.wikipedia.org/wiki/Bag-of-words_model) model can be used to create a feature vector that describes a document. It simply counts how many times each word occurs in a document, and uses that as a vector.

For example, the document "John likes to watch movies. Mary likes movies too." has the vector ```{"John":1,"likes":2,"to":1,"watch":1,"movies":2,"Mary":1,"too":1}```.

We will use a standardized vector of the ![most common 10,000 words](https://raw.githubusercontent.com/first20hours/google-10000-english/master/google-10000-english.txt) so that we can compare different documents and find similarities between them.

## Data 

For this example, I pulled 1000 books from ![Project Guttenberg](https://www.gutenberg.org/) for the dataset. I made a simple bash script ```books/get_books.sh``` that calls wget a bunch of times to download the books. Note there are two commonly used url formats, so I just try both for each book_id.

While getting a bag-of-words vector as quickly as possible might not be super important for books that are decades old, one could easily imagine applying this very same technique to tweets that are coming in real time, scraping news stories from major news orgs, or many other uses cases when speed is paramount.

## Parallel bag-of-words vector generation.

This problem is embarassingly parallel. Each book can be assigned to a goroutine that calculates its bag-of-words vector completely independently of every other book, so this makes a perfect use case for parallel programming and we should see some nice speedup!

## Term frequency-inverse document frequency 

![Term frequency-inverse document frequency (tf-idf)](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) takes this one step further. It is a slightly more sophisticated statistic that can be used to represent how important certain words are to a document. Instead of simply counting how many times a word appears, we include another term called *inverse document frequency*. It accounts for what percentage of our documents of *documents* it appears in as well. This way common words like 'the' 'an' don't get a big weight since they appear in every document, but rarer words get a bigger weight since they are more likely to be important.

To do this, we need to wait until all of the workers are finished counting their individual documents, and then we know how many documents each word appears in.

Once we know for each word how many documents it appears in, we can calculate the tfidf for each word. 