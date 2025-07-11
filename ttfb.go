package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func measureTTFBOverConn(ctx context.Context, conn net.Conn, host string) (ttfb time.Duration, err error) {
	if deadline, ok := ctx.Deadline(); ok {
		conn.SetDeadline(deadline)
		defer conn.SetDeadline(time.Time{})
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+host+"/", nil)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}
	req.Host = host
	req.Header.Set("Connection", "close")

	if err := req.Write(conn); err != nil {
		return 0, fmt.Errorf("write request: %w", err)
	}

	// Custom buffered reader to intercept first byte timing
	reader := bufio.NewReader(conn)

	// Use a 1-byte read to capture TTFB manually
	var firstByte [1]byte
	start := time.Now()
	_, err = reader.Read(firstByte[:])
	ttfb = time.Since(start)

	if err != nil {
		return ttfb, fmt.Errorf("error whilst reading first byte: %w", err)
	}

	return ttfb, nil
}
