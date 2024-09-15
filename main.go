package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/hashicorp/go-hclog"
	"gp-hip-report/internal/report"
	"syscall"
)

type CLI struct {
	Cookie     string `help:"Global Protect cookie"`
	MD5        string `help:"MD5 sum" name:"md5"`
	ClientIpv4 string `help:"Client IPv4 address" name:"client-ip"`
	ClientIpv6 string `help:"Client IPv4 address" name:"client-ipv6"`
	ClientOs   string `help:"Client OS (ignored)" name:"client-os"`
}

func setuidRoot() {
	err := syscall.Setuid(0)

	if err != nil {
		hclog.Default().Warn("failed to setuid", err)
	}
}

func (c *CLI) Run(kctx *kong.Context) error {
	ctx := context.Background()

	setuidRoot()

	hipReport, err := report.GenerateReport(ctx, c.Cookie, c.MD5, c.ClientIpv4, c.ClientIpv6)

	xmlReport, err := xml.MarshalIndent(hipReport, "", " ")

	if err != nil {
		return err
	}

	fmt.Println(xml.Header + string(xmlReport))

	return nil
}

func main() {
	cli := &CLI{}

	kctx := kong.Parse(cli)

	if err := kctx.Run(); err != nil {
		hclog.Default().Error("Failed to run command", err)
	}
}
