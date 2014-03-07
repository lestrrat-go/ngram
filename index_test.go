package ngram

import (
  "testing"
)

func buildNgramIndex(n int, inputs []string) *Index {
  i := NewIndex(n)
  for _, t := range inputs {
    i.AddString(t)
  }
  return i
}

func TestBasic(t *testing.T) {
  texts := []string {
    `abcdefghijklmnopqrstuvwxyz`,
    `1234567890abc1234567890`,
    `1234578990`,
  }
  i := buildNgramIndex(3, texts)

  matches := i.FindMatchingStrings(`abc def`)
  if len(matches) != 2 {
    t.Errorf("Expected 2 matches, got only %d", len(matches))
  } else {
    for _, s := range matches {
      if s.Content() != texts[0] && s.Content() != texts[1] {
        t.Errorf("Expected to NOT match %s", s)
      }
    }
  }

  // XXX This test sucks
  x := i.FindSimilarStrings(`abc def`)
  if len(x) != 2 {
    t.Errorf("Expected 2 items, got %d", len(x))
  }
}

func TestIndex_SimilarStrings(t *testing.T) {
  i := buildNgramIndex(3, []string { "abc", "abcabc", "aabc" });

  var ret []string

  ret = i.FindSimilarStrings("abc")
  if len(ret) != 3 {
    t.Logf("Expected to match 3 items, got %d", len(ret))
  }

  i.SetMinSimilarityScore(0.9)
  ret = i.FindSimilarStrings("abc")
  if len(ret) != 1 {
    t.Logf("Expected to match 1 item, got %d", len(ret))
  }
}