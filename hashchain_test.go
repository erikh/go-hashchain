package hashchain

import (
	"bytes"
	"crypto/sha512"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	const (
		expected1    = "6cdd012b8c0f1286f868ae4b9a7c6f559a5b83e1273b34eba9ac143dfa37097d37657aff9a492479ec310ae1af01bbbcea02676c52d1b9d6ee66481074343b9b"
		expected2    = "be3724850f3fbfd673052614631df78b53e882fab130d50792c9bdf97eb469679f355bcf43c0183ac94ecf0ab63bd7e40349357cc5932d1b60234887fe4631bc"
		sumExpected1 = "4665197f1f403a20406ce66c41a54144e813073b3c81d56a36a9c94a9506f459890a788bf384feda4966e76d1deb560135875544c26975f4a2faf8834912a144"
		sumExpected2 = "3655e269134d0d6c322eaf489492e842481e4bfbff5dd664903bf7d904791fd29038c44992f27160fccb05181926090231b542d88b71d10e285e9c7d9cdde941"
	)

	c := &Chain{}

	sum, err := c.Add(bytes.NewBuffer([]byte("buffer")), sha512.New())
	if err != nil {
		t.Fatal(err)
	}

	if sum != expected1 {
		t.Fatalf("Sum was not expected, was %q", sum)
	}

	sum, err = c.Sum(sha512.New())
	if err != nil {
		t.Fatal(err)
	}

	if sum != sumExpected1 {
		t.Fatalf("Sum was not expected, was %q", sum)
	}

	sum, err = c.Add(bytes.NewBuffer([]byte("buffer2")), sha512.New())
	if err != nil {
		t.Fatal(err)
	}

	if sum != expected2 {
		t.Fatalf("Sum was not expected, was %q", sum)
	}

	sum, err = c.Sum(sha512.New())
	if err != nil {
		t.Fatal(err)
	}

	if sum != sumExpected2 {
		t.Fatalf("Sum was not expected, was %q", sum)
	}

	if !reflect.DeepEqual(c.AllSums(), []string{expected1, expected2}) {
		t.Fatal("AllSums did not yield the correct result")
	}

	out := bytes.NewBuffer(nil)

	sum, err = c.AddInline(out, bytes.NewBuffer([]byte("buffer")), sha512.New())
	if err != nil {
		t.Fatal(err)
	}

	if sum != expected1 {
		t.Fatal("expected sum did not match")
	}

	if out.String() != "buffer" {
		t.Fatal("output was not yielded")
	}
}

func TestMatch(t *testing.T) {
	sums := []string{}

	chain1 := &Chain{}
	chain2 := &Chain{}

	same := []string{
		"buffer",
		"buffer2",
	}

	different1 := []string{
		"buffer1-1",
		"buffer1-2",
	}

	different2 := []string{
		"buffer2-1",
		"buffer2-2",
	}

	for _, str := range same {
		for _, chain := range []*Chain{chain1, chain2} {
			sum, err := chain.Add(bytes.NewBuffer([]byte(str)), sha512.New())
			if err != nil {
				t.Fatal(err)
			}

			sums = append(sums, sum)
		}
	}

	if chain, err := chain1.FirstMatch(chain2); err != nil || !reflect.DeepEqual(chain, chain2) {
		t.Fatal("FirstMatch did not equal chain2")
	}

	if chain, err := chain2.FirstMatch(chain1); err != nil || !reflect.DeepEqual(chain, chain1) {
		t.Fatal("FirstMatch did not equal chain1")
	}

	if chain, err := chain1.LastMatch(chain2); err != nil || !reflect.DeepEqual(chain, &Chain{}) {
		t.Fatal("LastMatch did not equal end of chain2")
	}

	if chain, err := chain2.LastMatch(chain1); err != nil || !reflect.DeepEqual(chain, &Chain{}) {
		t.Fatal("LastMatch did not equal end of chain1")
	}

	for _, str := range different1 {
		_, err := chain1.Add(bytes.NewBuffer([]byte(str)), sha512.New())
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, str := range different2 {
		_, err := chain2.Add(bytes.NewBuffer([]byte(str)), sha512.New())
		if err != nil {
			t.Fatal(err)
		}
	}

	if chain, err := chain1.FirstMatch(chain2); err != nil || !reflect.DeepEqual(chain, chain2) {
		t.Fatal("FirstMatch did not equal chain2")
	}

	if chain, err := chain2.FirstMatch(chain1); err != nil || !reflect.DeepEqual(chain, chain1) {
		t.Fatal("FirstMatch did not equal chain1")
	}

	if chain, err := chain1.LastMatch(chain2); err != nil || !reflect.DeepEqual(chain.chain, chain2.chain[2:]) {
		t.Fatal("LastMatch did not equal end of chain2")
	}

	if chain, err := chain2.LastMatch(chain1); err != nil || !reflect.DeepEqual(chain.chain, chain1.chain[2:]) {
		t.Fatal("LastMatch did not equal end of chain1")
	}

	chain1 = &Chain{}
	chain2 = &Chain{}

	for _, str := range different1 {
		_, err := chain1.Add(bytes.NewBuffer([]byte(str)), sha512.New())
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, str := range different2 {
		_, err := chain2.Add(bytes.NewBuffer([]byte(str)), sha512.New())
		if err != nil {
			t.Fatal(err)
		}
	}

	if _, err := chain1.FirstMatch(chain2); err == nil {
		t.Fatal("Chains matched and should not have")
	}

	if _, err := chain2.FirstMatch(chain1); err == nil {
		t.Fatal("Chains matched and should not have")
	}

	if _, err := chain1.LastMatch(chain2); err == nil {
		t.Fatal("Chains matched and should not have")
	}

	if _, err := chain2.LastMatch(chain1); err == nil {
		t.Fatal("Chains matched and should not have")
	}
}
