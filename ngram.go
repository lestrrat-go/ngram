package ngram

import(
  "unicode/utf8"
)

type Tokenizer interface {
  Tokenize()
  Tokens() []*Token
  Text() string
  N() int
}

type Tokenize struct {
  n int
  text string
  parsed bool
  tokens []*Token
}

type Token struct {
  n     *Tokenize
  start int
  end   int
}

func NewTokenize(n int, input string) *Tokenize {
  return &Tokenize {
    n,
    input,
    false,
    nil,
  }
}

func (n *Tokenize) N() int {
  return n.n
}

func (n *Tokenize) Text() string {
  return n.text
}

func (n *Tokenize) Parse() {
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

  // tokens hold the indices into the input string
  // there are ridx - 1 runes in this string. for an 'n'-gram,
  // there are ridx + 1 - ncount tokens
  ncount := n.N()
  tokens := make([]*Token, ridx + 1 - ncount)
  for i := 0; i <= ridx - ncount; i++ {
    end := i + ncount
    if end >= len(runes) {
      break
    }
    tokens[i] = n.NewToken(runes[i], runes[end])
  }
  n.tokens = tokens
  n.parsed = true
}

func (n *Tokenize) NewToken(start, end int) *Token {
  return &Token { n, start, end }
}

func (n *Tokenize) Tokens() []*Token {
  n.Parse()
  return n.tokens
}

func (s *Token) String() string {
  return s.n.text[s.start:s.end]
}

func (s *Token) Start() int {
  return s.start
}

func (s *Token) End() int {
  return s.end
}