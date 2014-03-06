package ngram

type Index struct {
  n int
  documents []string
  invertedIndex map[string][]int
}

func NewIndex(n int) *Index {
  return &Index { n, []string {}, map[string][]int {} }
}

func (i *Index) AddString(input string) {
  idx := len(i.documents)
  i.documents = append(i.documents, input)
  n := New(i.n, input)
  seen := make(map[string]bool)
  for _, s := range n.Segments() {
    str := s.String()
    _, ok := seen[str]
    if ok {
      continue
    }
    seen[str] = true

    list, ok := i.invertedIndex[str]
    if ! ok {
      i.invertedIndex[str] = []int {idx}
    } else {
      i.invertedIndex[str] = append(list, idx)
    }
  }
}

func (i *Index) FindSimilarString(input string) []string {
  return nil // TODO
}

func (i *Index) FindMatchingStrings(input string) []string {
  n := New(i.n, input)

  seen := make(map[int]bool)
  indices := []int {}
  for _, s := range n.Segments() {
    str := s.String()
    list, ok := i.invertedIndex[str]
    if ! ok {
      continue
    }

    for _, idx := range list {
      _, ok := seen[idx]
      if ok {
        continue
      }
      seen[idx] = true

      indices = append(indices, idx)
    }
  }

  ret := make([]string, len(indices))
  for idx, k := range indices {
    ret[idx] = i.documents[k]
  }
  return ret
}
