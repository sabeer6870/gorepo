package datastore

import (
	"fmt"

	"gopkg.in/couchbase/gocb.v1"
	"urllookupservice/common"
)

type Couchbase struct {
	Bucket *gocb.Bucket
 	Offline bool
}

type value struct {
	IsMalware bool `json:"IsMalware"`
}

func New() (*Couchbase, error){
	c := &Couchbase{
		Offline: true,
	}
	globalCluster, err := gocb.Connect("http://localhost:8091/")
	if err != nil{
		return c, err
	}

	globalCluster.Authenticate(gocb.PasswordAuthenticator{
		Username: common.CBUserName,
		Password: common.CBPassword,
	})

	globalBucket, err := globalCluster.OpenBucket(common.CouchbaseBucket, "")
	if err != nil{
		return c, err
	}
	c = &Couchbase{
		Bucket: globalBucket,
		Offline: false,
	}
	return c, nil
}

func (c *Couchbase) Get(key string) bool {
	fmt.Printf("Retrieving key:%s", key)
	var v value
	_, err := c.Bucket.Get(key, &v)
	if err != nil {
		return false
	}
	return v.IsMalware
}

func (c *Couchbase) PutOrPost(key string) {
	fmt.Printf("Adding key:%s\n", key)
	v := value{
		IsMalware: true,
	}
	c.Delete(key)
	c.Bucket.Insert(key, &v, common.Expiry)
}

func (c *Couchbase) Delete(key string) {
	fmt.Printf("Deleting key:%s\n", key)
	var v value
	cas, err := c.Bucket.Get(key, &v)
	if err == nil {
		c.Bucket.Remove(key, cas)
	}
}

// local test
func Test() {
	x, err := New()
	fmt.Println(err)
	fmt.Println(x.Get("google1"))
	x.PutOrPost("test1.com")
	fmt.Println(x.Get("test.com"))
	x.PutOrPost("google")
	x.Delete("google")
	x.PutOrPost("google.com")
}