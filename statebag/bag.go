package statebag

// Bag is a collection of state accessible by key. It is effectively a
// hash table that has slightly richer higher level functions built on top.
//
// Otto uses this throughout the core system in order to pass arbitrary
// data around.
//
// A Bag is not thread-safe. Concurrent read/write access should be protected
// by a lock. Concurrent read-only access is safe.
//
// NOTE: At time of writing, this is basically just syntactic sugar over
// a hash. Richer methods will be added over time.
type Bag struct {
	// Data is the data inside the state bag. This shouldn't be set
	// directly. The Get and Set methods should be used instead. Accessing
	// it directly can result in undefined behavior.
	Data map[string]interface{}
}

// Get reads data from the state bag
func (b *Bag) Get(k string) (interface{}, bool) {
	v, ok := b.Data[k]
	return v, ok
}

// Set sets data into the state bag
func (b *Bag) Set(k string, v interface{}) {
	b.Data[k] = v
}
