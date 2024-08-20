package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	randString := testTools.RandomString(10)
	if len(randString) != 10 {
		t.Error("wrong length of random string")
	}
}
