package provisioner

import (
	"bssms/internal/bssms"
	"context"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/lrstanley/go-bogon"
	_ "github.com/lrstanley/go-bogon"
	"net"
	"os"
)

func matchInterface(addr map[string]interface{}) (matched, internal bool, ip, mac string) {
	var tp interface{}
	var ok bool
	if tp, ok = addr["version"]; !ok {
		return
	}
	var ftp float64
	if ftp, ok = tp.(float64); !ok || ftp != 4 {
		return
	}
	// OS-EXT-IPS-MAC:mac_addr
	// OS-EXT-IPS:type fixed
	// addr 141.94.107.178
	sip4, ok := addr["addr"].(string)
	if !ok {
		return
	}
	if ipv4 := net.ParseIP(sip4); ipv4 == nil {
		return
	}
	tbogon, _ := bogon.Is(sip4)
	tmac, ok := addr["OS-EXT-IPS-MAC:mac_addr"].(string)
	if !ok {
		return
	}
	matched, internal, ip, mac = true, tbogon, sip4, tmac
	return
}

func GetOSServers() ([]InstallHost, error) {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}
	pc, err := openstack.AuthenticatedClient(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	cc, err := openstack.NewComputeV2(pc, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		return nil, err
	}
	listOpts := servers.ListOpts{}

	allPages, err := servers.List(cc, listOpts).AllPages(context.TODO())
	if err != nil {
		return nil, err
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		return nil, err
	}
	var ihs []InstallHost
	for _, server := range allServers {
		ih := InstallHost{
			Installable: bssms.Installable{
				Name:       server.Name,
				ServerUuid: server.ID,
			},
		}
		for _, ads := range server.Addresses {
			addrs := ads.([]interface{})
			for _, ad := range addrs {
				addr := ad.(map[string]interface{})
				matched, internal, ip, mac := matchInterface(addr)
				if matched {
					ih.MacAddress = mac
					if internal {
						ih.IPIntAddress = ip
					} else {
						ih.IPExtAddress = ip
					}
				}
			}
		}
		ihs = append(ihs, ih)
	}
	return ihs, nil
}
