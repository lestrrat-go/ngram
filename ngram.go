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

  // "runes" is not reaully an array of runes, it holds indices to start/end
  // of runes. so for example, if the string starts with a 3, 1, 2 byte runes,
  // we'd see runes = { 0, 3, 4, 6 ... }.
  //
  // To get the first trigram, we'd do input[runes[0]:runes[3]].
  // To get the first bigram, we'd do input[runes[0]:runes[2]], etc.
  //
  // runes is initialized as len(input) + 1 because we need to hold
  // the maximum number of runes (max is when all the input is 1 byte chars,
  // which is len(input)), plus the "end of string" marker. The end of string
  // marker is required because we always need start + end markers to access
  // the slice out of the original input. hence the + 1
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
  // there are ridx + 1 - ncount segments
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