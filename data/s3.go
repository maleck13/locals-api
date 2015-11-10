package data

import (
	"github.com/maleck13/locals-api/Godeps/_workspace/src/github.com/mitchellh/goamz/aws"
	"github.com/maleck13/locals-api/Godeps/_workspace/src/github.com/mitchellh/goamz/s3"
	"io/ioutil"
	"log"
	"os"
)

const (
	LOCALS_BUCKET = "localsie"
)

func PutInBucket(file string, remoteName string) (string, error) {
	var (
		_AWS_SECRET   string
		_AWS_ACCESS   string
		_AWS_LOCATION string
	)

	f, err := os.Open(file)

	_AWS_SECRET = os.Getenv("AWS_SECRET")
	_AWS_ACCESS = os.Getenv("AWS_ACCESS")
	_AWS_LOCATION = os.Getenv("AWS_LOCATION")

	if nil != err {
		return "", err
	}
	defer f.Close()
	auth, err := aws.GetAuth(_AWS_ACCESS, _AWS_SECRET)
	s3conn := s3.New(auth, aws.EUWest)
	bucket := s3conn.Bucket(LOCALS_BUCKET)
	data, err := ioutil.ReadAll(f)
	if nil != err {
		return "", err
	}
	log.Println(" adding filepath " + remoteName)

	err = bucket.Put(remoteName, data, "image/jpeg", s3.PublicRead)
	if nil != err {
		return "", err
	}
	return _AWS_LOCATION + "/" + remoteName, nil
}
