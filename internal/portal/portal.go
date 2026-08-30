package portal

import (
	"context"
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

type ScreenshotMode int

const (
	ScreenshotModeFull      = iota
	ScreenshotModeSelection = iota
)

const (
	portalPrefix         = "org.freedesktop.portal"
	dbusObjDestination   = portalPrefix + ".Desktop"
	dbusObjPath          = "/org/freedesktop/portal/desktop"
	dbusScreenshotMethod = portalPrefix + ".Screenshot.Screenshot"
	dbusReqInterface     = portalPrefix + ".Request"
)

type screenshotResponse struct {
	ResponseCode uint32
	URI          string
}

type Coordinate struct {
	X1, X2, Y1, Y2 int
}

type ScreenshotOptions struct {
	Mode      ScreenshotMode
	Selection Coordinate
}

type Portal interface {
	Screenshot(context.Context, ScreenshotOptions) (string, error)
}

type Client struct {
	conn *dbus.Conn
}

func parseBody(s dbus.Signal) screenshotResponse {
	respCode := s.Body[0].(uint32)
	res := s.Body[1].(map[string]dbus.Variant)
	var uri string

	uriVariant, ok := res["uri"]

	if ok {
		uri = uriVariant.Value().(string)
	}

	return screenshotResponse{
		ResponseCode: respCode,
		URI:          uri,
	}
}

func New() (*Client, error) {
	conn, err := dbus.ConnectSessionBus()

	if err != nil {
		return nil, fmt.Errorf("error creating DBUS connection: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Screenshot(ctx context.Context, clientOpts ScreenshotOptions) (string, error) {
	var options map[string]dbus.Variant
	if clientOpts.Mode == ScreenshotModeFull {
		options = map[string]dbus.Variant{
			"interactive": dbus.MakeVariant(false),
		}
	} else {
		return "", errors.New("only Fullscreen mode is enabled at this moment")
	}

	obj := c.conn.Object(
		dbusObjDestination,
		dbus.ObjectPath(dbusObjPath),
	)

	var requestPath dbus.ObjectPath

	err := obj.Call(
		dbusScreenshotMethod,
		0,
		"",
		options,
	).Store(&requestPath)

	if err != nil {
		return "", fmt.Errorf("failed to call DBUS screenshot method: %w", err)
	}

	err = c.conn.AddMatchSignal(
		dbus.WithMatchObjectPath(requestPath),
		dbus.WithMatchInterface(dbusReqInterface),
		dbus.WithMatchMember("Response"),
	)

	if err != nil {
		return "", fmt.Errorf("failed to add signal match on DBUS request: %w", err)
	}

	signals := make(chan *dbus.Signal, 10)

	c.conn.Signal(signals)
	defer c.conn.RemoveSignal(signals)

	for signal := range signals {
		if signal.Path != requestPath {
			continue
		}

		if signal.Name != portalPrefix+".Request.Response" {
			continue
		}

		body := parseBody(*signal)

		if body.ResponseCode != 0 {
			return "", errors.New("couldn't get the screenshot")
		}

		if body.URI == "" {
			return "", errors.New("couldn't get a non-empty URI")
		}

		return body.URI, nil
	}

	return "", errors.New("failed to get screenshot for an unknown reason")
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close DBUS connection: %w", err)
	}

	return nil
}
