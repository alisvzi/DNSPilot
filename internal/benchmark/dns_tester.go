package benchmark

import (
	"context"
	"time"

	"github.com/miekg/dns"
)

func TestDNS(ip string) Result {

	start := time.Now()

	client := dns.Client{
		Timeout: 2 * time.Second,
	}

	msg := dns.Msg{}
	msg.SetQuestion("google.com.", dns.TypeA)

	_, _, err := client.ExchangeContext(
		context.Background(),
		&msg,
		ip+":53",
	)

	if err != nil {

		return Result{
			ServerIP: ip,
			Success:  false,
		}
	}

	return Result{
		ServerIP: ip,
		Latency:  time.Since(start).Milliseconds(),
		Success:  true,
	}
}
