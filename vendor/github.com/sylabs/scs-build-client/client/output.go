// Copyright (c) 2018-2022, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// GetOutput streams build output for the provided buildID to w. The context controls the lifetime
// of the request.
func (c *Client) GetOutput(ctx context.Context, buildID string, w io.Writer) error {
	u := c.baseURL.ResolveReference(&url.URL{
		Path: "v1/build-ws/" + buildID,
	})

	wsScheme := "ws"
	if c.baseURL.Scheme == "https" {
		wsScheme = "wss"
	}
	u.Scheme = wsScheme

	h := http.Header{}
	c.setRequestHeaders(h)

	dialer := websocket.DefaultDialer

	// Due to this issue (https://github.com/gorilla/websocket/issues/601), it is not possible
	// clone the 'c.HTTPClient' transport, so we take only the InsecureSkipVerify and RootCAs
	// parameters.
	if tr, ok := c.httpClient.Transport.(*http.Transport); ok {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: tr.TLSClientConfig.InsecureSkipVerify,
			RootCAs:            tr.TLSClientConfig.RootCAs,
		}
	}

	ws, resp, err := dialer.DialContext(ctx, u.String(), h)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer resp.Body.Close()
	defer ws.Close()

	errChan := make(chan error)

	go func() {
		defer close(errChan)
		errChan <- func() error {
			for {
				// Read from websocket
				mt, r, err := ws.NextReader()
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return nil
				} else if err != nil {
					return fmt.Errorf("failed to read output: %w", err)
				}

				if mt != websocket.TextMessage {
					continue
				}

				if _, err := io.Copy(w, r); err != nil {
					return fmt.Errorf("failed to copy output: %w", err)
				}
			}
		}()
	}()

	select {
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = c.Cancel(ctx, buildID) //nolint:contextcheck

		ws.Close()

		<-errChan
		return nil
	case err := <-errChan:
		return err
	}
}
