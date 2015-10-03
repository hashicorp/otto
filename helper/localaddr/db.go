package localaddr

import (
	"bytes"
	"container/heap"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

var (
	boltLocalAddrBucket = []byte("localaddr")
	boltBuckets         = [][]byte{
		boltLocalAddrBucket,
	}

	boltVersionKey  = []byte("version")
	boltAddrMapKey  = []byte("addr_map")
	boltAddrHeapKey = []byte("addr_heap")
	boltSubnetKey   = []byte("subnet")
)

var (
	boltDataVersion byte = 2
)

var boltCidr *net.IPNet

func init() {
	_, cidr, err := net.ParseCIDR("100.64.0.0/10")
	if err != nil {
		panic(err)
	}

	boltCidr = cidr
}

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
		data := bucket.Get(boltSubnetKey)
		if data == nil {
			panic("no subnet")
		}

		// Get the existing IP addresses that we've mapped
		addrMap, addrQ, err := this.getData(bucket)
		if err != nil {
			return err
		}

		// Parse the subnet
		ip, ipnet, err := net.ParseCIDR(string(data))
		if err != nil {
			return err
		}
		ip = ip.To4()

		// Generate a random IP in our subnet and try to use it
		found := false
		for {
			// Create a random address in the subnet
			ipRaw := make([]byte, 4)
			binary.LittleEndian.PutUint32(ipRaw, rand.Uint32())
			for i, v := range ipRaw {
				ip[i] = ip[i] + (v &^ ipnet.Mask[i])
			}

			// If this IP exists, then try again
			if _, ok := addrMap[ip.String()]; ok {
				continue
			}

			// We found an IP!
			found = true
			break
		}

		// If we didn't find an IP, we just use the oldest one available
		if !found {
			result = heap.Pop(&addrQ).(*ipEntry).Value
		}

		// Set the result
		result = ip

		// Add the IP to the map
		entry := &ipEntry{LeaseTime: time.Now().UTC(), Value: ip}
		heap.Push(&addrQ, entry)
		addrMap[ip.String()] = entry.Index

		// Store the data
		return this.putData(bucket, addrMap, addrQ)
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
		bucket := tx.Bucket(boltLocalAddrBucket)

		// Get the existing IP addresses that we've mapped
		addrMap, addrQ, err := this.getData(bucket)
		if err != nil {
			return err
		}

		// If it isn't in there, we're done
		idx, ok := addrMap[ip.String()]
		if !ok {
			return nil
		}

		// Delete and save
		delete(addrMap, ip.String())
		addrQ, addrQ[len(addrQ)-1] = append(addrQ[:idx], addrQ[idx+1:]...), nil
		heap.Init(&addrQ)

		return this.putData(bucket, addrMap, addrQ)
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
		bucket := tx.Bucket(boltLocalAddrBucket)

		// Get the existing IP addresses that we've mapped
		addrMap, addrQ, err := this.getData(bucket)
		if err != nil {
			return err
		}

		// If it isn't in there, we're done
		idx, ok := addrMap[ip.String()]
		if !ok {
			return nil
		}

		entry := addrQ[idx]
		entry.LeaseTime = time.Now().UTC()
		addrQ.Update(entry)

		return nil
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
	bootstrap := false
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLocalAddrBucket)
		data := bucket.Get([]byte("version"))
		if data == nil || len(data) == 0 {
			version = boltDataVersion
			bootstrap = true
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

	// Map of update functions
	updateMap := map[byte]func(*bolt.DB) error{
		1: this.v1_to_v2,
	}
	for version < boltDataVersion {
		log.Printf(
			"[INFO] upgrading lease DB from v%d to v%d", version, version+1)
		err := updateMap[version](db)
		if err != nil {
			return nil, fmt.Errorf(
				"Error upgrading data from v%d to v%d: %s",
				version, version+1, err)
		}

		version++
	}

	// Bootstrap if we have to
	if bootstrap {
		if err := this.v1_to_v2(db); err != nil {
			return nil, err
		}
	}

	// Just call the upgrade to init
	return db, nil
}

func (this *DB) v1_to_v2(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLocalAddrBucket)

		// Replace the subnet with the fixed CIDR we're using. The old CIDR
		// doesn't matter...
		err := bucket.Put(boltSubnetKey, []byte(boltCidr.String()))
		if err != nil {
			return err
		}

		boltAddrBucket := []byte("addrs")
		if b := bucket.Bucket(boltAddrBucket); b != nil {
			// And delete all the addresses, we start over
			err = bucket.DeleteBucket(boltAddrBucket)
			if err != nil {
				return err
			}
		}

		return bucket.Put(boltVersionKey, []byte{byte(2)})
	})
}

func (this *DB) putData(
	bucket *bolt.Bucket,
	addrMap map[string]int,
	addrQ ipQueue) error {
	var buf, buf2 bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(addrMap); err != nil {
		return err
	}
	if err := bucket.Put(boltAddrMapKey, buf.Bytes()); err != nil {
		return err
	}

	if err := gob.NewEncoder(&buf2).Encode(addrQ); err != nil {
		return err
	}
	if err := bucket.Put(boltAddrHeapKey, buf2.Bytes()); err != nil {
		return err
	}

	return nil
}

func (this *DB) getData(bucket *bolt.Bucket) (map[string]int, ipQueue, error) {
	var addrQ ipQueue
	heapRaw := bucket.Get(boltAddrHeapKey)
	if heapRaw == nil {
		addrQ = ipQueue(make([]*ipEntry, 0, 1))
	} else {
		dec := gob.NewDecoder(bytes.NewReader(heapRaw))
		if err := dec.Decode(&addrQ); err != nil {
			return nil, nil, err
		}
	}

	var addrMap map[string]int
	mapRaw := bucket.Get(boltAddrMapKey)
	if mapRaw == nil {
		addrMap = make(map[string]int)
	} else {
		dec := gob.NewDecoder(bytes.NewReader(mapRaw))
		if err := dec.Decode(&addrMap); err != nil {
			return nil, nil, err
		}
	}

	return addrMap, addrQ, nil
}
