package bolt

import (
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/boltdb/bolt"
	"golang.org/x/net/context"

	"github.com/gnhuy91/profilesvc"
)

// NewBoltService construct a new BoltDB Service and implements
// profilesvc.Service.
func NewBoltService(db *bolt.DB, bucketName string) profilesvc.Service {
	return &Service{DB: db, bucketName: bucketName}
}

// Service represents a BoltDB implementation of profilesvc.Service.
type Service struct {
	DB         *bolt.DB
	bucketName string
}

func (s *Service) PostProfile(ctx context.Context, p profilesvc.Profile) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte(s.bucketName))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		p.ID = strconv.Itoa(int(id))

		// Marshal user data into bytes.
		buf, err := json.Marshal(p)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put([]byte(p.ID), buf)
	})
}

func (s *Service) GetProfile(ctx context.Context, id string) (profilesvc.Profile, error) {
	var profile profilesvc.Profile
	var err error

	err = s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		profileBytes := b.Get([]byte(id))

		err := json.Unmarshal(profileBytes, &profile)
		if err != nil {
			return err
		}
		return nil
	})

	return profile, err
}

func (s *Service) DeleteProfile(ctx context.Context, id string) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(s.bucketName))
		return b.Delete([]byte(id))
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
