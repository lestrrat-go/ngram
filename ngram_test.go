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
    n := NewTokenize(i, input)
    tokens := n.Tokens()
    if len(tokens) != runeCount + 1 - i {
      t.Errorf("There should be %d tokens, but got %d", runeCount + 1 - i, len(tokens))
    }

    for _, s := range tokens {
      count := utf8.RuneCountInString(s.String())
      if count != i {
        t.Errorf("Expected segment to have %d runes, got %d", i, count)
      }
    }
  }
}

func TestTokenize_TokenSet(t *testing.T) {
  input := `the quick red quick red fox red fox jumps fox jumps over jumps over the over the lazy the lazy brown lazy brown dog`
  tok := NewTokenize(3, input)
  set := tok.TokenSet()

  // Check uniqueness
  uniq := map[string]bool {}
  for str := range set.Iter() {
    if _, ok := uniq[str.(string)]; ok {
      t.Errorf("token %s appears multiple times!", str)
    }
    uniq[str.(string)] = true
  }
}

func ExampleTokenize() {
  input := `the quick red quick red fox red fox jumps fox jumps over jumps over the over the lazy the lazy brown lazy brown dog`
  n := NewTokenize(3, input) // Trigram
  for _, s := range n.Tokens() {
    log.Printf("segment = %s", s)
  }
}
