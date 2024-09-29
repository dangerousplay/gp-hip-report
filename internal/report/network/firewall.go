package network

import (
	"encoding/json"
	"encoding/xml"
	"github.com/cockroachdb/errors"
	"github.com/hashicorp/go-hclog"
	"github.com/life4/genesis/slices"
	"gp-hip-report/internal/utils"
	"os/exec"
	"strings"
)

const (
	firewallEntryName = "firewall"

	ufwFirewall      = "ufw"
	iptablesFirewall = "iptables"
	nftablesFirewall = "nftables"
)

var knownFirewallVendors = map[string]string{
	"ufw":      "Canonical Ltd.",
	"iptables": "IPTables",
	"nftables": "The Netfilter Project",
}

type Firewall struct {
	XMLName xml.Name        `xml:"entry"`
	Name    string          `xml:"name,attr"`
	List    []FirewallEntry `xml:"list>entry"`
}

type FirewallEntry struct {
	ProductInfo FirewallProductInfo `xml:"ProductInfo"`
}

type FirewallProductInfo struct {
	Prod      FirewallProd `xml:"Prod"`
	IsEnabled string       `xml:"is-enabled"`
}

type FirewallProd struct {
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
	Vendor  string `xml:"vendor,attr"`
}

type nftables struct {
	Nftables []nftable `json:"nftables"`
}

type nftable struct {
	Metainfo struct {
		Version           string `json:"version"`
		ReleaseName       string `json:"release_name"`
		JsonSchemaVersion int    `json:"json_schema_version"`
	} `json:"metainfo,omitempty"`
	Table struct {
		Family string `json:"family"`
		Name   string `json:"name"`
		Handle int    `json:"handle"`
	} `json:"table,omitempty"`
	Chain struct {
		Family string `json:"family"`
		Table  string `json:"table"`
		Name   string `json:"name"`
		Handle int    `json:"handle"`
		Type   string `json:"type,omitempty"`
		Hook   string `json:"hook,omitempty"`
		Prio   int    `json:"prio,omitempty"`
		Policy string `json:"policy,omitempty"`
	} `json:"chain,omitempty"`
	Rule struct {
		Family string `json:"family"`
		Table  string `json:"table"`
		Chain  string `json:"chain"`
		Handle int    `json:"handle"`
		Expr   []struct {
			Match struct {
				Op   string `json:"op"`
				Left struct {
					Meta struct {
						Key string `json:"key"`
					} `json:"meta,omitempty"`
					Payload struct {
						Protocol string `json:"protocol"`
						Field    string `json:"field"`
					} `json:"payload,omitempty"`
				} `json:"left"`
				Right interface{} `json:"right"`
			} `json:"match,omitempty"`
			Counter struct {
				Packets int `json:"packets"`
				Bytes   int `json:"bytes"`
			} `json:"counter,omitempty"`
			Accept interface{} `json:"accept"`
			Xt     struct {
				Type string `json:"type"`
				Name string `json:"name"`
			} `json:"xt,omitempty"`
			Jump struct {
				Target string `json:"target"`
			} `json:"jump,omitempty"`
			Drop  interface{} `json:"drop"`
			Limit struct {
				Rate  int    `json:"rate"`
				Burst int    `json:"burst"`
				Per   string `json:"per"`
			} `json:"limit,omitempty"`
			Return interface{} `json:"return"`
		} `json:"expr"`
	} `json:"rule,omitempty"`
}

func GetFirewallInfo() (*Firewall, error) {
	firewall := &Firewall{
		Name: firewallEntryName,
		List: []FirewallEntry{},
	}

	logger := hclog.Default().Named("firewall")

	var errs error

	ufw, err := checkUfw()

	if err != nil {
		logger.Warn("Failed to check if ufw is active", "err", err)
		errs = errors.CombineErrors(errs, err)
	}

	if ufw != nil {
		firewall.List = append(firewall.List, *ufw)
	}

	nft, err := checkNft()

	if err != nil {
		logger.Warn("Failed to check if nftables is active", "err", err)
		errs = errors.CombineErrors(errs, err)
	}

	if nft != nil {
		firewall.List = append(firewall.List, *nft)
	}

	iptables, err := checkIptables()

	if err != nil {
		logger.Warn("Failed to check if iptables is active", "err", err)
		errs = errors.CombineErrors(errs, err)
	}

	if iptables != nil {
		firewall.List = append(firewall.List, *iptables)
	}

	return firewall, errs
}

func checkUfw() (*FirewallEntry, error) {
	ufwInstalled, ufwVersion, err := utils.CheckExists("ufw", []string{"ufw", "--version"})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !ufwInstalled {
		return nil, nil
	}

	version := parseUfwVersion(ufwVersion)
	active, err := isUfwEnabled()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	entry := &FirewallEntry{ProductInfo: FirewallProductInfo{
		Prod: FirewallProd{
			Name:    ufwFirewall,
			Version: version,
			Vendor:  knownFirewallVendors[ufwFirewall],
		},
		IsEnabled: utils.BoolToString(active),
	}}

	return entry, nil
}

func checkNft() (*FirewallEntry, error) {
	installed, versionOutput, err := utils.CheckExists("nft", []string{"nft", "--version"})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !installed {
		return nil, nil
	}

	version := parseNftVersion(versionOutput)

	entry := &FirewallEntry{
		ProductInfo: FirewallProductInfo{
			Prod: FirewallProd{
				Name:    nftablesFirewall,
				Version: version,
				Vendor:  knownFirewallVendors[nftablesFirewall],
			},
		},
	}

	cmd := exec.Command("nft", "-j", "list", "ruleset")
	output, err := cmd.Output()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	var rules nftables

	err = json.Unmarshal(output, &rules)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	active := true

	if len(rules.Nftables) < 1 {
		active = false
	}

	active = slices.Any(rules.Nftables, func(el nftable) bool {
		return len(el.Rule.Expr) > 0
	})

	entry.ProductInfo.IsEnabled = utils.BoolToString(active)

	return entry, nil
}

func checkIptables() (*FirewallEntry, error) {
	ipTablesInstalled, ipTablesVersion, err := utils.CheckExists("iptables", []string{"iptables", "--version"})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !ipTablesInstalled {
		return nil, nil
	}

	version := parseIptablesVersion(ipTablesVersion)

	entry := &FirewallEntry{ProductInfo: FirewallProductInfo{
		Prod: FirewallProd{
			Name:    iptablesFirewall,
			Version: version,
			Vendor:  knownFirewallVendors[iptablesFirewall],
		},
		IsEnabled: utils.BoolToString(false),
	}}

	return entry, nil
}

func parseNftVersion(output string) string {
	parts := strings.Split(output, " ")

	if len(parts) < 2 {
		return ""
	}

	return strings.TrimPrefix(parts[1], "v")
}

func parseIptablesVersion(output string) string {
	parts := strings.Split(output, " ")

	if len(parts) < 2 {
		return ""
	}

	return strings.TrimPrefix(parts[1], "v")
}

func parseUfwVersion(output string) string {
	parts := strings.Split(strings.Split(output, "\n")[0], " ")

	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func isUfwEnabled() (bool, error) {
	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()

	if err != nil {
		return false, err
	}

	lines := strings.Split(string(output), "\n")

	if len(lines) < 1 {
		return false, nil
	}

	statusParts := strings.Split(lines[0], ":")
	if len(statusParts) < 2 {
		return false, nil
	}

	status := strings.TrimSpace(statusParts[1])
	active := strings.EqualFold(status, "active")

	return active, nil
}
