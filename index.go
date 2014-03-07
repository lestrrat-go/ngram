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

func (i *Index) GetMinScore() float64 {
  return i.minSimilarityScore
}

func (i *Index) SetMinScore(min float64) {
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

type MatchResult struct {
  Score float64
  Item  IndexItem
}

// search for similar strings in index, sending the search results
// to the given channel
func (i *Index) IterateSimilar(input string, min float64) <-chan MatchResult {
  c := make(chan MatchResult)
  go i.iterateSimilar(c, input, min)
  return c
}

func (i *Index) iterateSimilar(c chan MatchResult, input string, min float64) {
  n := NewTokenize(i.n, input)
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
      if score >= min {
        c <- MatchResult { score, item.item }
      }
    }
  }
  close(c)
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

func (i *Index) FindSimilarStrings(input string) []string {
  c := i.IterateSimilar(input, i.minSimilarityScore)

  ret := []string {}
  for r := range c {
    ret = append(ret, r.Item.Content())
  }
  return ret
}

func (i *Index) FindBestMatch(input string) string {
  c := i.IterateSimilar(input, i.minSimilarityScore)

  maxScore := 0.0
  var best IndexItem
  for r := range c {
    if maxScore < r.Score {
      maxScore = r.Score
      best = r.Item
    }
  }
  return best.Content()
}

func (i *Index) FindMatchingStrings(input string) []string {
  c := i.IterateSimilar(input, 0)

  ret := []string {}
  for r := range c {
    ret = append(ret, r.Item.Content())
  }
  return ret
}

