# gp-hip-report

**Global Protect HIP Report** generation as a replacement for PanGpHip, currently for Linux. 

<!-- TOC -->
* [gp-hip-report](#gp-hip-report)
  * [How to use](#how-to-use)
    * [Using with gpclient](#using-with-gpclient)
    * [Using alone](#using-alone)
  * [Install](#install)
    * [Debian](#debian)
    * [Manual](#manual)
      * [Linux](#linux)
  * [Supported Operating Systems](#supported-operating-systems)
  * [Supported features](#supported-features)
    * [Disk encryption](#disk-encryption)
    * [Network](#network)
      * [Firewall](#firewall)
    * [Anti Malware](#anti-malware)
<!-- TOC -->

## How to use

[![asciicast](https://asciinema.org/a/gXwmIzQefDh8K0NmZBNOkyRAX.svg)](https://asciinema.org/a/gXwmIzQefDh8K0NmZBNOkyRAX)

You can use it with `gpclient`, `openconnect` or alone.

```text
Usage: gp-hip-report [flags]

Flags:
  -h, --help                  Show context-sensitive help.
      --cookie=STRING         Global Protect cookie
      --md5=STRING            MD5 sum
      --client-ip=STRING      Client IPv4 address
      --client-ipv6=STRING    Client IPv4 address
      --client-os=STRING      Client OS (ignored)
```

### Using with gpclient

```shell
gpclient connect --csd-wrapper "/usr/bin/gp-hip-report" <PORTAL_URL>
```

### Using alone

Generating a HIP report without MD5, Client IPv4/IPv6, username, Domain information:
```shell
gp-hip-report
```

Generating a HIP report with user, domain, computer, Client IP information:
```shell
gp-hip-report --md5 3e33a6232b6a99c625d5e2696492c714 --cookie "user=dangerousplay&domain=net&computer=wts02" --client-ip 192.168.1.2
```

## Install

### Debian

Download the `.deb` from a release and installs it:
```shell
wget https://github.com/dangerousplay/gp-hip-report/releases/download/<version>/gp-hip-report_<version>_<arch>.deb
dpkg -i gp-hip-report_<version>_<arch>.deb
```

### Manual

#### Linux

Download the precompiled binary and add the `setuid` flag.
```shell
$ sudo wget https://github.com/dangerousplay/gp-hip-report/releases/download/<version>/gp-hip-report_<version>_linux_<arch> -O /usr/bin/gp-hip-report
$ sudo chmod u+s /usr/bin/gp-hip-report
```

## Supported Operating Systems

| OS     | Versions |
|--------|----------|
| Ubuntu | 24.04    |

## Supported features

### Disk encryption

Reports the encryption state of the mounted disks on the system.

| Name       | OS    |
|------------|-------|
| cryptsetup | Linux |


### Network

Reports the interfaces Mac Address, assigned IPv4/IPv6 addresses and the status of installed firewalls.

#### Firewall

| Name     | OS    |
|----------|-------|
| iptables | Linux |
| nftables | Linux |
| ufw      | Linux |


### Anti Malware

| Name          | OS    |
|---------------|-------|
| Falcon Sensor | Linux |


### Patch management

| OS    | Name | Missing patch report? |
|-------|------|-----------------------|
| Linux | apt  | No                    |


