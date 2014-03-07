package ngram

import (
  "crypto/md5"
  "errors"
  "fmt"
)

// ngram.Index is a naive implementation of an in-memmory index using
// ngram tokens as its keys. It allows one to register strings, which
// are tokenized and registered in an inverted index to the strings
// they originated from, allowing us to search from a ngram token to
// strings containing them.
type ngrameKey struct {
  id string // id of IndexItem
  ntoken int // number of ngram tokens in the document
}

type InvertedIndex map[string]map[string]int
type Index struct {
  n int
  items IndexItemDepot
  // maps ngram key to indexitemwithmetadata's key
  // a map of item.id => int is used so that we don't store
  // redundant ids
  invertedIndex InvertedIndex
}

type IndexItem interface {
  Id() string
  Content() string
}

type Document struct {
  id string
  content string
}

// Map of item id to IndexItemWithMetadata.
// This object holds the IndexItem itself, as well as the 
// number of ngrams in this document
type IndexItemDepot map[string]*IndexItemWithMetadata
type IndexItemWithMetadata struct {
  item IndexItem
  ngrams int
}

func NewDocument(id, content string) *Document {
  if id == "" {
    h := md5.New()
    h.Write([]byte(content))
    id = fmt.Sprintf("%x", h.Sum(nil))
  }
  return &Document { id, content }
}

func (d *Document) Id() string {
  return d.id
}

func (d *Document) Content() string {
  return d.content
}

func NewIndex(n int) *Index {
  return &Index { n, IndexItemDepot {} , InvertedIndex {} }
}

func (i *Index) AddString(input string) (error) {
  return i.AddItem(NewDocument("", input))
}

func (i *Index) AddItem(item IndexItem) (error) {
  // XXX mutex?

  id := item.Id()
  if _, ok := i.items[id]; ok {
    return errors.New(
      fmt.Sprintf(
        "Item %s already exists in index",
        id,
      ),
    )
  }
  seen := make(map[string]bool)

  n := NewTokenize(i.n, item.Content())
  tokens := n.Tokens()
  for _, s := range tokens {
    str := s.String()
    _, ok := seen[str]
    if ok {
      continue
    }
    seen[str] = true

    h, ok := i.invertedIndex[str]
    if ! ok {
      h = map[string]int {}
      i.invertedIndex[str] = h
    }
    h[id]++
  }

  ntokens := len(tokens)
  i.items[id] = &IndexItemWithMetadata {
    item,
    ntokens,
  }
  return nil
}

func (i *Index) FindSimilarString(input string) []string {
  return nil // TODO
}

func (i *Index) FindMatchingStrings(input string) []IndexItem {
  n := NewTokenize(i.n, input)

  seen := make(map[string]bool)
  indices := []string {}
  for _, s := range n.Tokens() {
    str := s.String()
    h, ok := i.invertedIndex[str]
    if ! ok {
      continue
    }

    for key, _ := range h {
      _, ok := seen[key]
      if ok {
        continue
      }
      seen[key] = true

      indices = append(indices, key)
    }
  }

  ret := make([]IndexItem, len(indices))
  for idx, k := range indices {
    ret[idx] = i.items[k].item
  }
  return ret
}
