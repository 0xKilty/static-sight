**Host-based signatures** - indicators that certains files have been impacted by the malware, help figure what the malware has done to the system, not the malware itself

**Network-based signatures** - indicators that tell what the malware has done on a network, similar to host-based signatures

### Types of analysis
**Basic Statis Analysis** - Examination of the file with out looking at the instructions, confirmation if a file is malicious, and can help create network signatures

**Basic Dynamic Analysis** - Running the malware and observing its behavior. Gain methods to remove the malware, and generate effective signatures

**Advanced Static Analysis** - Disassembling the malware to view the source code

**Advanced Dynamic Analysis** - Figuring out how the disasembled code works and looks for detailed information on the malware

### Types of Malware
* Backdoor - allows a hacker to gain access onto a computer
* Botnet - Similar to backdoor, but the infected computers are given the same instructions from a command and control server
* Downloader - malware that downloads other programs, often more malware. Typically just malware that gains access to a system
* Information stealing - Malware that collects information for a computer. Ex. Sniffers, password hash grabbers, and keyloggers. 
* Launcher - starts other malware, usually to launch malware such that the launched malware is stealthy and undetected
* Root kit - Malware to hide other malware by disguising or hiding
* Scareware - malware to scare the user into doing something, posing as antimalware that the user should buy only to remove the scareware
* Spam-sending - malware that allows hackers to use the infected computer to send spam
* Worm or virus - Malware that copies itself and infects other computers

* Mass malware - malware not meant for a specific target (less complex)
* Targeted malware - malware specialized for a certain target (more complex)

### General Rules
* Don't try to understand every detail, try to figure out how the malware works and what the malware does
* Don't get stuck trying one approach, there are many ways and tools to analyze malware
* It's a cat and mouse game, malware is a game with opponents leading to new discoveries all the time

