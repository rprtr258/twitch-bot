package balaboba

import "testing"

func TestGenerate(t *testing.T) {
	c := ClientRus

	gen, err := c.Generate("123", Standart)
	if err != nil {
		t.Fatal(err)
	}
	if gen.BadQuery {
		t.Fatal("bad query", gen.BadQuery)
	}
}
