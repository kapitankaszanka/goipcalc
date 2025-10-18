# GoIPCalc

A simple IP subnet calculator written in Go.  
You provide an IPv4/IPv6 address in CIDR notation (e.g. `192.168.1.42/27`, `::2001:db8:1/89`), and the program calculates basic network information.

> ⚠️ Disclaimer: I’m currently learning Go, and this is my **first project** in this language — so treat it as a learning experiment rather than a production-ready tool. 😊

## 📦 What this tool does

Given an IP address (usually in `IP/MASK` CIDR format), it will:

- ✅ Validate the input address `custom validation`
- ✅ Calculate the **network address**  
- ✅ Calculate the **broadcast address**  
- ✅ Display the **subnet mask** (CIDR and dotted-decimal)  
- ✅ Show the **first and last usable host addresses**  
- ✅ Calculate the **total number of hosts** in the subnet  

---

## ⚙️ Requirements

- Go (recommended: `1.20+`)

---

## 🚀 Usage

Run directly without building:

```bash
go run . 192.168.1.42/27
```
Or build a binary:
```
go build -o goipcalc
./goipcalc 192.168.1.42/27
```

### basic usage
```
goipcalc --help
Usage: goipcalc [OPTIONS] [ADDR/PLEN]
Examples:
  goipcalc -d 10.0.0.1/24
  goipcalc 2001:db8::1/64 192.168.10.11/28
Options:
  [ADDR/PLEN] address/prefix lenght, can be multiple
  -d    IPv4 address to calculate
  -j    json output
  -json-indent
        change json output to indentation
```
```
goipcalc 2001:db8::1/64 192.168.1.24/25
---
Full address:  2001:db8:0:0:0:0:0:1/64
Network:       2001:db8:0:0:0:0:0:0
Last address:  2001:db8:0:0:ffff:ffff:ffff:ffff
---
Full address:  192.168.1.24/25
Network:       192.168.1.0
Broadcast:     192.168.1.127
```
```
goipcalc -d 2001:db8::1/64 192.168.1.24/25
---
Full address:  2001:db8:0:0:0:0:0:1/64
Network:       2001:db8:0:0:0:0:0:0
Last address:  2001:db8:0:0:ffff:ffff:ffff:ffff
Address:       2001:db8:0:0:0:0:0:1
Mask:          64
Mask address:  ffff:ffff:ffff:ffff:0:0:0:0
Hosts number:  18 446 744 073 709 551 616
---
Full address:  192.168.1.24/25
Network:       192.168.1.0
Broadcast:     192.168.1.127
Address:       192.168.1.24
Mask:          25
Mask address:  255.255.255.128
Hosts number:  128
```
```
goipcalc -j 2001:db8::1/64 192.168.1.24/25 555.555.555.555/24
{"results":[{"full_address":"2001:db8:0:0:0:0:0:1/64","network":"2001:db8:0:0:0:0:0:0","last_address":"2001:db8:0:0:ffff:ffff:ffff:ffff"},{"full_address":"192.168.1.24/25","network":"192.168.1.0","broadcast":"192.168.1.127"}],"errors":["skip \"555.555.555.555/24\": invalid address: 555.555.555.555\n"]}
```
```
goipcalc -j -json-indent 2001:db8::1/64 192.168.1.24/25
{
  "results": [
    {
      "full_address": "2001:db8:0:0:0:0:0:1/64",
      "network": "2001:db8:0:0:0:0:0:0",
      "last_address": "2001:db8:0:0:ffff:ffff:ffff:ffff"
    },
    {
      "full_address": "192.168.1.24/25",
      "network": "192.168.1.0",
      "broadcast": "192.168.1.127"
    }
  ]
}
```
```
goipcalc -d -j -json-indent 2001:db8::1/64 192.168.1.24/25
{
  "results": [
    {
      "full_address": "2001:db8:0:0:0:0:0:1/64",
      "network": "2001:db8:0:0:0:0:0:0",
      "last_address": "2001:db8:0:0:ffff:ffff:ffff:ffff",
      "address": "2001:db8:0:0:0:0:0:1",
      "mask": 64,
      "mask_address": "ffff:ffff:ffff:ffff:0:0:0:0",
      "hosts_number": 18446744073709551616
    },
    {
      "full_address": "192.168.1.24/25",
      "network": "192.168.1.0",
      "broadcast": "192.168.1.127",
      "address": "192.168.1.24",
      "mask": 25,
      "mask_address": "255.255.255.128",
      "hosts_number": 128
    }
  ]
}
```

