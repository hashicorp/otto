package localaddr

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

var (
	boltLocalAddrBucket = []byte("localaddr")
	boltAddrBucket      = []byte("addrs")
	boltBuckets         = [][]byte{
		boltLocalAddrBucket,
	}
)

var (
	boltDataVersion byte = 1
)

// DB is a database of local addresses, and provides operations to find
// the next available address, release an address, etc.
//
// DB will act as an LRU: if there are no available IP addresses, it will find
// the oldest IP address and give that to you. This is to combat the fact that
// the things that use IP addresses can often be killed outside of our control,
// and the oldest one is most likely to be stale. This should be an edge
// case.
//
// The first time DB is used, it will find a usable subnet space and
// allocate that as its own. After it allocates that space, it will use
// that for the duration of this DBs existence. The usable subnet space
// is randomized to try to make it unlikely to have a collision.
//
// DB uses a /24 so the entire space of available IP addresses is only
// 256, but these IPs are meant to be local, so they shouldn't overflow
// (it would mean more than 256 VMs are up... or that each of those VMs
// has a lot of network interfaces. Both cases are unlikely in Otto).
//
// FUTURE TODO:
//
//   * Allocate additional subnets once we run out of IP space (vs. LRU)
//
type DB struct {
	// Path is the path to the IP database. This file doesn't need to
	// exist but needs to be a writable path. The parent directory will
	// be made.
	Path string
}

// Next returns the next IP that is not allocated.
func (this *DB) Next() (net.IP, error) {
	db, err := this.db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var result net.IP
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLocalAddrBucket)
		data := bucket.Get([]byte("subnet"))
		if data == nil {
			panic("no subnet")
		}

		// Get the bucket with addresses
		addrBucket, err := bucket.CreateBucketIfNotExists(boltAddrBucket)
		if err != nil {
			return err
		}

		// Parse the subnet
		ip, _, err := net.ParseCIDR(string(data))
		if err != nil {
			return err
		}

		// Start the IP at a random number in the range. We add 2 to
		// avoid ".1" which is usually used by the gateway.
		var start byte = byte(rand.Int31n(253) + 2)
		ipKey := len(ip) - 1
		ip[ipKey] = start

		// Increment the IP until we find one we don't have
		var oldestIP net.IP
		var oldestTime string
		for {
			key := []byte(ip.String())
			data := addrBucket.Get(key)
			if data != nil {
				// We can just use a lexical comparison of time because
				// the formatting allows it.
				if dataStr := string(data); oldestTime == "" || dataStr < oldestTime {
					oldestIP = ip
					oldestTime = dataStr
				}

				// Increment the IP
				ip[ipKey]++
				if ip[ipKey] == 0 {
					ip[ipKey] = 2
				}
				if ip[ipKey] != start {
					continue
				}

				// Return the oldest one if we're out
				ip = oldestIP
				key = []byte(ip.String())
			}

			// Found one! Insert it and return it
			err := addrBucket.Put(key, []byte(time.Now().UTC().String()))
			if err != nil {
				return err
			}

			result = ip
			return nil
		}
	})

	return result, err
}

// Release releases the given IP, removing it from the database.
func (this *DB) Release(ip net.IP) error {
	db, err := this.db()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLocalAddrBucket).Bucket(boltAddrBucket)
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(ip.String()))
	})
}

// Renew updates the last used time of the given IP to right now.
//
// This should be called whenever a DB-given IP is used to make sure
// it isn't chosen as the LRU if we run out of IPs.
func (this *DB) Renew(ip net.IP) error {
	db, err := this.db()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLocalAddrBucket).Bucket(boltAddrBucket)
		if bucket == nil {
			return nil
		}

		key := []byte(ip.String())
		return bucket.Put(key, []byte(time.Now().UTC().String()))
	})
}

// db returns the database handle, and sets up the DB if it has never
// been created.
func (this *DB) db() (*bolt.DB, error) {
	// Make the directory to store our DB
	if err := os.MkdirAll(filepath.Dir(this.Path), 0755); err != nil {
		return nil, err
	}

	// Create/Open the DB
	db, err := bolt.Open(this.Path, 0644, nil)
	if err != nil {
		return nil, err
	}

	// Create the buckets
	err = db.Update(func(tx *bolt.Tx) error {
		for _, b := range boltBuckets {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Check the DB version
	var version byte
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLocalAddrBucket)
		data := bucket.Get([]byte("version"))
		if data == nil || len(data) == 0 {
			version = boltDataVersion
			return bucket.Put([]byte("version"), []byte{boltDataVersion})
		}

		version = data[0]
		return nil
	})
	if err != nil {
		return nil, err
	}

	if version > boltDataVersion {
		return nil, fmt.Errorf(
			"IP data version is higher than this version of Otto knows how\n"+
				"to handle! This version of Otto can read up to version %d,\n"+
				"but version %d data file found.\n\n"+
				"This means that a newer version of Otto touched this data,\n"+
				"or the data was corrupted in some other way.",
			boltDataVersion, version)
	}

	// Init the subnet
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLocalAddrBucket)
		data := bucket.Get([]byte("subnet"))

		// If we already have a subnet, bail
		if data != nil {
			return nil
		}

		// No subnet, allocate one and save it
		ipnet, err := UsableSubnet()
		if err != nil {
			return err
		}

		return bucket.Put([]byte("subnet"), []byte(ipnet.String()))
	})

	return db, nil
}
