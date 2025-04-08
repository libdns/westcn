package westcn

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/libdns/libdns"
)

func TestProvider(t *testing.T) {
	provider := Provider{
		Username:    "",
		APIPassword: "",
	}

	fmt.Println(provider.GetRecords(context.Background(), ""))

	fmt.Println(provider.AppendRecords(context.Background(), "", []libdns.Record{
		libdns.RR{
			Name: "sub",
			TTL:  10 * time.Minute,
			Type: "A",
			Data: "8.8.8.8",
		},
	}))

	fmt.Println(provider.SetRecords(context.Background(), "", []libdns.Record{
		libdns.RR{
			Name: "sub",
			TTL:  10 * time.Minute,
			Type: "A",
			Data: "8.8.8.8",
		},
	}))

	fmt.Println(provider.DeleteRecords(context.Background(), "", []libdns.Record{
		libdns.RR{
			Name: "sub",
			TTL:  10 * time.Minute,
			Type: "A",
			Data: "8.8.8.8",
		},
	}))
}
