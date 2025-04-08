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

	recs, err := provider.GetRecords(context.Background(), "")
	if err != nil {
		t.Fatalf("failed to get records: %v", err)
	}
	fmt.Println(recs)

	recs, err = provider.AppendRecords(context.Background(), "", []libdns.Record{
		libdns.RR{
			Name: "sub",
			TTL:  10 * time.Minute,
			Type: "A",
			Data: "8.8.8.8",
		},
	})
	if err != nil {
		t.Fatalf("failed to append records: %v", err)
	}
	fmt.Println(recs)

	recs, err = provider.SetRecords(context.Background(), "", []libdns.Record{
		libdns.RR{
			Name: "sub",
			TTL:  10 * time.Minute,
			Type: "A",
			Data: "8.8.8.8",
		},
	})
	if err != nil {
		t.Fatalf("failed to set records: %v", err)
	}
	fmt.Println(recs)

	recs, err = provider.DeleteRecords(context.Background(), "", []libdns.Record{
		libdns.RR{
			Name: "sub",
			TTL:  10 * time.Minute,
			Type: "A",
			Data: "8.8.8.8",
		},
	})
	if err != nil {
		t.Fatalf("failed to delete records: %v", err)
	}
	fmt.Println(recs)
}
