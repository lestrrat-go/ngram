package ngram

import (
  "testing"
)

func TestSimilarText(t *testing.T) {
  i := NewIndex(3)
  texts := []string {
    `abcdefghijklmnopqrstuvwxyz`,
    `1234567890abc1234567890`,
    `1234578990`,
  }
  for _, t := range texts {
    i.AddString(t)
  }

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
}