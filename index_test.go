package ngram

import (
  "log"
  "testing"
)

func ExampleIndex () {
  trigram := NewIndex(3)
  inputs  := []string {
    `...`,
    `...`,
    `...`,
  }
  for _, v := range inputs {
    trigram.AddString(v)
  }

  // Find strings whose scores are above trigram.GetMinScore()
  // (which is by default 0)
  matches := trigram.FindSimilarStrings(`...`)
  log.Printf("%#v", matches)

  // Find 1 best match (the best score) out of similar strings
  best := trigram.FindBestMatch(`...`)
  log.Printf("%s", best)

  // Iterate match results
  minScore := 0.5
  limit := 0
  c := trigram.IterateSimilar(` ... your input ...`, minScore, limit)
  for r := range c {
    log.Printf("Item id %s matched with score %d", r.Item.Id(), r.Score)
    log.Printf("Content of Item was %s", r.Item.Content())
  }
}

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
      if s != texts[0] && s != texts[1] {
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

  i.SetMinScore(0.9)
  ret = i.FindSimilarStrings("abc")
  if len(ret) != 1 {
    t.Logf("Expected to match 1 item, got %d", len(ret))
  }
}

func TestIndex_FindBestMatch(t *testing.T) {
  i := buildNgramIndex(3, []string { "abc", "abcabc", "aabc" })
  best := i.FindBestMatch("abc")
  if best != "abc" {
    t.Errorf("Expected 'abc', got '%s'", best)
  }
}