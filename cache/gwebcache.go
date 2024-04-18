package cache

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

type void struct{}
type Vendor string

var sources = map[Vendor][]string{
	"bazooka": {"http://cache.trillinux.org/g2/bazooka.php"},
}

func QueryBazookaVendors(ctx context.Context) []string {
	allIPs := make(map[string]void)

	for _, u := range sources["bazooka"] {
		uri, _ := url.Parse(u)
		q := uri.Query()
		q.Set("showhosts", "1")
		q.Set("net", "gnutella2")
		uri.RawQuery = q.Encode()
		ips, _ := getBazookaIPs(ctx, uri)

		for _, ip := range ips {
			allIPs[ip] = void{}
		}
	}

	var allIPsSlice []string
	for ip := range allIPs {
		allIPsSlice = append(allIPsSlice, ip)
	}
	return allIPsSlice
}

func getBazookaIPs(ctx context.Context, uri *url.URL) ([]string, error) {
	const filter = `gnutella:host:(\d*.\d*.\d*.\d*:\d*)`
	ipsRegexp := regexp.MustCompile(filter)

	req, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		fmt.Println(err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ips := ipsRegexp.FindAllString(string(body), -1)

	return ips, nil
}
