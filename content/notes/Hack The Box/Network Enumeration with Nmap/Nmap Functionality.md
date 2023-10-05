### Scan multiple IPs
- 10.129.2.0/24
- 10.129.2.18 10.129.2.19 10.129.2.20
- 10.129.2.18-20
- ``` nmap -iL hosts.lst ```

### Flags
- ``` -sn ``` disables port scanning
- ``` -oA <file> ``` creates files that store the results of the scan in many formats
- ``` -iL <file> ``` scans all the IPs specified in the file
- ``` -PE ``` performs the ping scan by using 'ICMP Echo requests' against the target
- ``` --packet-trace ``` Shows all packets sent and received
- 