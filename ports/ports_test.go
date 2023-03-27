package ports

import "testing"

var nextExpects = []int{1,2,3,4}

func TestNext(t *testing.T) {

  next := Next()

  for _, expected := range nextExpects {
    if output := next(); output != expected {
      t.Errorf("Output %q not equal to expected %q", output, expected )
    }
  }
}
