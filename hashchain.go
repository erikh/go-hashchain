package hashchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
)

var (
	ErrNoMatch = errors.New("Chains do not match")
)

// Chain is a chain of sums iconifying generational changes to a single file,
// so that you can trace the origin of divergence later.
// Add to the Chain with the Add() function.
type Chain struct {
	chain [][]byte
}

// Add a file's sum to the chain. Returns the hex encoded sum for convenience.
// Always use the same hash.Hash type, but never the same object, with this
// function.
func (c *Chain) Add(r io.Reader, h hash.Hash) (string, error) {
	if _, err := io.Copy(h, r); err != nil {
		return "", fmt.Errorf("Could not digest reader: %w", err)
	}

	sum := h.Sum(nil)
	c.chain = append(c.chain, sum)

	return hex.EncodeToString(sum), nil
}

// Same as Add(), but also write it to the writer in the process. The sum is
// returned along with any error.
func (c *Chain) AddInline(w io.Writer, r io.Reader, h hash.Hash) (string, error) {
	reader := io.TeeReader(r, h)

	if _, err := io.Copy(w, reader); err != nil {
		return "", fmt.Errorf("Could not copy from reader: %w", err)
	}

	sum := h.Sum(nil)
	c.chain = append(c.chain, sum)

	return hex.EncodeToString(sum), nil
}

// Obtain the list of sums this chain possesses. Useful for marshaling.
func (c *Chain) AllSums() []string {
	res := []string{}

	for _, byt := range c.chain {
		res = append(res, hex.EncodeToString(byt))
	}

	return res
}

// Find the first match in the provided chain; returns the passed chain that
// follows the match.
func (c *Chain) FirstMatch(c2 *Chain) (*Chain, error) {
	for _, cbyt := range c.chain {
		for x, c2byt := range c2.chain {
			if bytes.Equal(c2byt, cbyt) {
				return &Chain{chain: c2.chain[x:]}, nil
			}
		}
	}

	return nil, ErrNoMatch
}

// Find the last match in the provided chain; returns the passed chain that
// follows the match, which may be empty.
func (c *Chain) LastMatch(c2 *Chain) (*Chain, error) {
	everEqual := false
	for _, cbyt := range c.chain {
		idx := 0
		equal := false

		for x, c2byt := range c2.chain {
			if bytes.Equal(c2byt, cbyt) {
				equal = true
				everEqual = true
				break
			}

			idx = x
		}

		if !equal && everEqual {
			if idx == 0 {
				return c2, nil
			}

			return &Chain{chain: c2.chain[idx-1:]}, nil
		}
	}

	if !everEqual {
		return nil, ErrNoMatch
	}

	return &Chain{}, nil // all match, so return an empty chain
}

// Sum the entire chain's sums. Useful for quickly determining a mismatch.
func (c *Chain) Sum(h hash.Hash) (string, error) {
	for _, byt := range c.chain {
		if _, err := io.Copy(h, bytes.NewBuffer(byt)); err != nil {
			return "", fmt.Errorf("Failed to hash chain: %w", err)
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
