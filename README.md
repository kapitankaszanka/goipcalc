# IPCalculator

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
- ✅ (Optionally) display the **wildcard mask** or **binary representation**

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
go build -o ipcalculator
./ipcalculator 192.168.1.42/27
```

### basic usage
```
ipcalculator -ip4 10.1.11.24/14 -ip 8.13.29.11/29 -ip6 2001:db8::/56 -ip6 2001:db::13:31:faf1/34

------ Start
--- Version IPv4
Addr/Pref     : 10.1.11.24/14
Address       : 10.1.11.24
Mask          : 255.252.0.0
Network       : 10.0.0.0
Broadcast     : 10.3.255.255
Host number   : 262142
--- Version IPv4
Addr/Pref     : 8.13.29.11/29
Address       : 8.13.29.11
Mask          : 255.255.255.248
Network       : 8.13.29.8
Broadcast     : 8.13.29.15
Host number   : 6
--- Version IPv6
Addr/Pref      : 2001:db8:0:0:0:0:0:0/56
Address        : 2001:db8:0:0:0:0:0:0
Mask           : 56
Network        : 2001:db8:0:0:0:0:0:0/56
Last address   : 2001:db8:0:ff:ffff:ffff:ffff:ffff/56
Host number    : To many to bother...
--- Version IPv6
Addr/Pref      : 2001:db:0:0:0:13:31:faf1/34
Address        : 2001:db:0:0:0:13:31:faf1
Mask           : 34
Network        : 2001:db:0:0:0:0:0:0/34
Last address   : 2001:db:3fff:ffff:ffff:ffff:ffff:ffff/34
Host number    : To many to bother...
------ End

```

