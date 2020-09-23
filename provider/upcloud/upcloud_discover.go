// Package upcloud provides node discovery for UpCloud.
package upcloud

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

type Provider struct {
	svc *service.Service
}

func (p *Provider) Help() string {
	return `upcloud:

    provider:       "upcloud"
    username: 	    The UpCloud account name
	password:       The UpCloud password
	zone:           The UpCloud zone to filter on
	tag:            The tag to filter on
	address_access: "utility", "public". (default: "utility")
	address_family: "IPv4", "IPv6" (default: "IPv4")
`
}

func (p *Provider) listServersByTagZone(tag string, zone string, l *log.Logger) ([]upcloud.Server, error) {
	res := []upcloud.Server{}
	servers, err := p.svc.GetServers()
	if err != nil {
		return nil, fmt.Errorf("error getting server list: %v", err)
	}
	for _, server := range servers.Servers {
		if server.Zone != zone {
			continue
		}
		for _, t := range server.Tags {
			if t == tag {
				res = append(res, server)
			}
		}
	}
	l.Printf("[DEBUG] discover-upcloud: Found %d servers with tag '%s' on zone '%s'", len(res), tag, zone)
	return res, nil
}

func (p *Provider) fetchAddrsFromServer(server *upcloud.Server, access string, family string, l *log.Logger) ([]string, error) {
	var addrs []string
	details, err := p.svc.GetServerDetails(&request.GetServerDetailsRequest{
		UUID: server.UUID,
	})
	if err != nil {
		return nil, fmt.Errorf("discover-upcloud: Fetching details for server %v failed: %s", server.UUID, err)
	}
	for _, addr := range details.IPAddresses {
		if addr.Access == access && addr.Family == family {
			addrs = append(addrs, addr.Address)
		}
	}
	return addrs, nil
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "upcloud" {
		return nil, fmt.Errorf("discover-upcloud: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	username := args["username"]
	password := args["password"]
	zone := args["zone"]
	tag := args["tag"]
	addressAccess := args["address_access"]
	addressFamily := args["address_family"]

	if addressAccess == "" {
		addressAccess = "utility"
	}
	if addressFamily == "" {
		addressFamily = "IPv4"
	}

	l.Printf("[INFO] discover-upcloud: Username is %q", username)

	p.svc = service.New(client.New(username, password))
	_, err := p.svc.GetAccount()
	if err != nil {
		return nil, fmt.Errorf("discover-upcloud: %s", err)
	}

	servers, err := p.listServersByTagZone(tag, zone, l)
	if err != nil {
		return nil, fmt.Errorf("error getting server list: %v", err)
	}

	l.Printf("[INFO] discover-upcloud: Found %d servers with tag '%s'", len(servers), tag)

	var addrs []string
	for _, server := range servers {
		if server.Zone == zone {
			l.Printf("[DEBUG] Instance UUID: %q", server.UUID)
			newAddrs, err := p.fetchAddrsFromServer(&server, addressAccess, addressFamily, l)
			if err != nil {
				return nil, err
			}
			addrs = append(addrs, newAddrs...)
		}
	}
	return addrs, nil
}
