package ngram

import (
  "log"
  "unicode/utf8"
  "testing"
)

func TestNgramASCII(t *testing.T) {
  input := `the quick red quick red fox red fox jumps fox jumps over jumps over the over the lazy the lazy brown lazy brown dog`
  testNgrams(t, input)
}

func TestNgramJapanese(t *testing.T) {
  input := `日本語でいろんな文章を書いてみよう`
  testNgrams(t, input)
}

func testNgrams(t *testing.T, input string) {
  runeCount := utf8.RuneCountInString(input)
  for i := 3; i < runeCount; i++ {
    n := New(i, input)
    segments := n.Segments()
    if len(segments) != runeCount + 1 - i {
      t.Errorf("There should be %d segments, but got %d", runeCount + 1 - i, len(segments))
    }

    for _, s := range segments {
      count := utf8.RuneCountInString(s.String())
      if count != i {
        t.Errorf("Expected segment to have %d runes, got %d", i, count)
      }
    }
  }
}

func ExampleNgram() {
  input := `the quick red quick red fox red fox jumps fox jumps over jumps over the over the lazy the lazy brown lazy brown dog`
  n := New(3, input) // Trigram
  for _, s := range n.Segments() {
    log.Printf("segment = %s", s)
  }
}
