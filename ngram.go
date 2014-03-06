package ngram

import(
  "unicode/utf8"
)

type NgramParser interface {
  Text() string
  N() int
}

type Ngram struct {
  n int
  text string
  parsed bool
  segments []*ngramSegment
}

type ngramSegment struct {
  n     *Ngram
  start int
  end   int
}

func New(n int, input string) *Ngram {
  return &Ngram {
    n,
    input,
    false,
    nil,
  }
}

func (n *Ngram) N() int {
  return n.n
}

func (n *Ngram) Text() string {
  return n.text
}

func (n *Ngram) Parse() {
  if n.parsed {
    return
  }

  input := n.text
  runes := make([]int, len(input) + 1)
  ridx  := 0 // rune index
  bidx  := 0 // byte index
  for len(input) > 0 {
    _, width := utf8.DecodeRuneInString(input)
    input = input[width:]
    runes[ridx] = bidx
    bidx += width
    ridx++
  }
  runes[ridx] = len(n.text)

  // segments hold the indices into the input string
  // there are ridx - 1 runes in this string. for an 'n'-gram,
  // there are maximum of ridx + 1 - ncount
  ncount := n.N()
  segments := make([]*ngramSegment, ridx + 1 - ncount)
  for i := 0; i <= ridx - ncount; i++ {
    end := i + ncount
    if end >= len(runes) {
      break
    }
    segments[i] = n.newSegment(runes[i], runes[end])
  }
  n.segments = segments
  n.parsed = true
}

func (n *Ngram) newSegment(start, end int) *ngramSegment {
  return &ngramSegment { n, start, end }
}

func (n *Ngram) Segments() []*ngramSegment {
  n.Parse()
  return n.segments
}

func (s *ngramSegment) String() string {
  return s.n.text[s.start:s.end]
}

func (s *ngramSegment) Start() int {
  return s.start
}

func (s *ngramSegment) End() int {
  return s.end
}