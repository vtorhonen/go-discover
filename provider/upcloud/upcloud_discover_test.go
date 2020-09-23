package upcloud_test

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/upcloud"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":       "upcloud",
		"username":       os.Getenv("UPCLOUD_USERNAME"),
		"password":       os.Getenv("UPCLOUD_PASSWORD"),
		"tag":            "go-discover-test-tag",
		"zone":           "de-fra1",
		"address_access": "utility",
		"address_family": "IPv4",
	}

	if args["username"] == "" || args["password"] == "" {
		t.Skip("Upcloud credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &upcloud.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
