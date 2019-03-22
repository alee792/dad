package dad

import (
	"bufio"
	"context"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Chain parses and stores word level n-grams
// and generates Markov chains.
type Chain struct {
	Config Config
	// Key: n-gram, value: possible suffixes.
	grams map[string][]string
	// Use of sync.Map is deferred to avoid excessive type assertions.
	mux *sync.Mutex
}

// Config for a Chain
type Config struct {
	// Order is the token count of an n-gram.
	Order int
}

// NewChain uses a Config to return a Chain
// with sensible defaults.
func NewChain(cfg Config) *Chain {
	if cfg.Order < 1 {
		cfg.Order = 2
	}
	rand.Seed(time.Now().Unix())
	c := &Chain{
		grams:  make(map[string][]string),
		mux:    &sync.Mutex{},
		Config: cfg,
	}
	return c
}

// Read corpus to generate pairs using Config of the Chain.
func (c *Chain) Read(ctx context.Context, r io.Reader) {
	scn := bufio.NewScanner(r)
	for scn.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			c.ReadSentence(scn.Text())
		}
	}
}

// ReadSentence parses a string of zero or more space-delimited
// sentences and incorporates them into the chain.
// See https://golang.org/doc/codewalk/markov/ for prior art.
// This implementation trades elegance for simplicity.
func (c *Chain) ReadSentence(s string) {
	words := strings.Fields(strings.TrimSpace(s))
	order := c.Config.Order
	lastIndex := len(words) - 1

	for i := range words {
		var gram string
		var next []string
		nextStart, nextEnd := i+order, i+(order*2)

		// Avoid out of index exception for orphans.
		if nextStart > lastIndex {
			gram = strings.Join(words[i:], " ")
			c.putGram(gram, nil)
			return
		}

		gram = strings.Join(words[i:nextStart], " ")

		// Check the end as well.
		if nextEnd > lastIndex {
			next = words[nextStart:]
		} else {
			next = words[nextStart:nextEnd]
		}

		// Compact any >1-grams into a single string.
		c.appendGram(gram, strings.Join(next, " "))
	}
}

// Generate a Dad joke using the chain's current corpus.
func (c *Chain) Generate(ctx context.Context) string {
	return c.GenerateWithMax(ctx, 20)
}

// GenerateWithMax words.
func (c *Chain) GenerateWithMax(ctx context.Context, max int) string {
	var grams []string
	gram := c.randomGram()
	for i := 0; i < max; i++ {
		// End of sentence.
		sfxs := c.getGram(gram)
		if len(sfxs) == 0 {
			break
		}
		next := sfxs[rand.Intn(len(sfxs))]
		grams = append(grams, next)
		gram = next
	}
	return strings.Join(grams, " ")
}

// The methods below are thread-safe, mutex-locked
// accessors for `Chain.gram`.

// Grams returns the n-grams and suffixes of the Markov Chain.
func (c *Chain) Grams() map[string][]string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.grams
}

func (c *Chain) getGram(k string) []string {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.grams[k]
}

func (c *Chain) randomGram() string {
	c.mux.Lock()
	defer c.mux.Unlock()
	var g string
	for g = range c.grams {
		break
	}
	return g
}

func (c *Chain) putGram(k string, v []string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	// Avoid overwriting gram values that might seem terminal!
	if v == nil {
		if _, ok := c.grams[k]; ok {
			return
		}
		v = []string{}
	}
	c.grams[k] = append(c.grams[k], v...)
}

// appendGram is a convenient way to append a single value.
func (c *Chain) appendGram(k, v string) {
	c.putGram(k, []string{v})
}
