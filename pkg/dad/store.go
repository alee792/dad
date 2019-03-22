package dad

// Storage defines how a corpus and n-grams are persisted.
// TODO: Shape up this API to be more generally applicable.
type Storage interface {
	Save(path string, grams map[string][]string) error
	Load(path string) (map[string][]string, error)
}

// Save NGrams to store.
func (c *Chain) Save(path string) error {
	return c.Store.Save(path, c.grams)
}

// Load NGrams from store.
func (c *Chain) Load(path string) (int, error) {
	grams, err := c.Store.Load(path)
	if err != nil {
		return 0, err
	}
	for k, v := range grams {
		c.putGram(k, v)
	}
	return len(grams), nil
}
