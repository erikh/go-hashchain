## go-hashchain: Verify a file changed the same way over generations

hashchain is not blockchain or a merkle tree, it is something much simpler. I
needed to track whether a file changed over a time the same way over many
generations. The verifiability of the "chain" is not as important as the
changes that the chain consists of.

Thus, this library was born. It's crypto built on somewhat solid foundations
but it's still _my_ crypto and I am _not_ a cryptographer. I don't suggest you
use it without evaluating what you need it for with a very critical eye.

### Usage

```go
// let's make two chains so we can show some things.
c1 := &Chain{}
c2 := &Chain{}

// the sum is yielded for the current provided buffer as a convenience for
// third-party tracking. Any hash.Hash interface works as a digest algorithm.

// we'll sum the same thing twice at first
sum, err := c1.Add(bytes.NewBuffer([]byte("buffer")), sha512.New())
if err != nil {
  log.Fatal(err)
}

sum, err = c2.Add(bytes.NewBuffer([]byte("buffer")), sha512.New())
if err != nil {
  log.Fatal(err)
}

// and one more time, for posterity's sake
sum, err = c1.Add(bytes.NewBuffer([]byte("buffer2")), sha512.New())
if err != nil {
  log.Fatal(err)
}

sum, err = c2.Add(bytes.NewBuffer([]byte("buffer2")), sha512.New())
if err != nil {
  log.Fatal(err)
}

// and now, we'll generate a mismatch
sum, err = c1.Add(bytes.NewBuffer([]byte("buffer1-1")), sha512.New())
if err != nil {
  log.Fatal(err)
}

sum, err = c2.Add(bytes.NewBuffer([]byte("buffer1-2")), sha512.New())
if err != nil {
  log.Fatal(err)
}

// difference functions. if the chains don't match at all, you'll get an error

// will return a chain containing c2's hashes, starting right after the first hash
chain, err := c1.FirstMatch(c2)

// will return a chain containing c2's hashes, starting right after the last hash
chain, err = c1.LastMatch(c2)

// we can tell the chains are different by comparing the sum of the entire chain:

fmt.Println(c1.Sum(sha512.New()) == c2.Sum(sha512.New()))
```

### Author

Erik Hollensbe <erik+github@hollensbe.org>

### License

MIT
