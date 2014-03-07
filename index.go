package ngram

import (
  "crypto/md5"
  "errors"
  "fmt"
  "math"
  "github.com/deckarep/golang-set"
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

  minSimilarityScore float64
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
  // ngram.String() => count
  ngrams mapset.Set
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
  return &Index { n, IndexItemDepot {} , InvertedIndex {}, 0.0 }
}

func (i *Index) SetMinSimilarityScore(min float64) {
  i.minSimilarityScore = min
}

func (i *Index) GetItemWithMetadata(id string) *IndexItemWithMetadata {
  return i.items[id]
}

func (i *Index) GetItem(id string) IndexItem {
  return i.items[id].item
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

  n := NewTokenize(i.n, item.Content())
  tokens := n.Tokens()
  set := mapset.NewSet()
  for _, s := range tokens {
    str := s.String()
    if set.Contains(str) {
      continue
    }

    set.Add(str)
    h, ok := i.invertedIndex[str]
    if ! ok {
      h = map[string]int {}
      i.invertedIndex[str] = h
    }
    h[id]++
  }

  i.items[id] = &IndexItemWithMetadata {
    item,
    set,
  }
  return nil
}

func (i *Index) FindSimilarStrings(input string) []string {
  items := i.FindSimilarItems(input)
  ret   := make([]string, len(items))
  for k, item := range items {
    ret[k] = item.Content()
  }
  return ret
}

func (i *Index) FindSimilarItems(input string) []IndexItem {
  n := NewTokenize(i.n, input)

  ret := []IndexItem {}
  inputset := n.TokenSet()
  seen := mapset.NewSet()
  for s := range inputset.Iter() {
    // for each token, find matching document
    itemids, ok := i.invertedIndex[s.(string)]
    if !ok {
      continue
    }

    for itemid := range itemids {
      if seen.Contains(itemid) {
        continue
      }

      seen.Add(itemid)

      item := i.GetItemWithMetadata(itemid)
      score := i.computeSimilarity(inputset, item.ngrams)
      if score >= i.minSimilarityScore {
        ret = append(ret, item.item)
      }
    }
  }

  return ret
}

func (i *Index) computeSimilarity(inputset, target mapset.Set) float64 {
  intersection := inputset.Intersect(target)
  total := target.Cardinality()
  contained := intersection.Cardinality()
  diff := total - contained

  totalExp := math.Pow(float64(total), 1.0)
  diffExp  := math.Pow(float64(diff), 1.0)

  return (totalExp - diffExp) / totalExp
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
