package ports

import "testing"

func TestNext(t *testing.T) {
  var nextExpects = []int{3000,3001,3002,3003,3004,3005,60010,60011,60012,60013,60014,60015}

  next := InitPortIterator("3000-3005,60010-60015")

  for _, expected := range nextExpects {
    if output, _ := next(); output != expected {
      t.Errorf("Output %q not equal to expected %q", output, expected )
    }
  }
}

func TestNextError(t *testing.T) {
  next := InitPortIterator("3000-3005")

  for i := 0; i < 10; i++ {
    _, err := next()
    if err != nil && err.Error() != "no more ports left" {
      t.Errorf("Error is not as expected: %q", err)
    }
  }
}
