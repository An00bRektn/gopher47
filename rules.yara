rule Gopher47_Linux {
    meta:
        description = "Gopher47 - file gopher47"
        author = "An00bRektn <An00bRektn@proton.me>"
        date = "2023-02-23"
    strings:
        $x1 = "cb701b6f0a2f55e3c269f5dde3f4ba25f55be6e65add1657b6843430bf1a4940" ascii

        $s1 = "X-IsAGopher" ascii
        $s2 = "[!] Cannot use execute-assembly in Linux!" ascii
        $s3 = "{\"task\":\"register\",\"data\":" ascii
        $s4 = "[!] Not currently supported for Linux!" ascii
    condition:
        uint16(0) == 0x457f and filesize > 5MB and
        1 of ($x*) or 3 of ($s*)
}

rule Gopher47_Windows {
    meta:
        description = "Gopher47 - file gopher47"
        author = "An00bRektn <An00bRektn@proton.me>"
        date = "2023-02-23"
    strings:
        $x1 = "cb701b6f0a2f55e3c269f5dde3f4ba25f55be6e65add1657b6843430bf1a4940" ascii

        $s1 = "X-IsAGopher" ascii
        $s2 = "[!] Failed to load assembly:" ascii
        $s3 = "{\"task\":\"register\",\"data\":" ascii
        $s4 = "[!] Failed to change protections on memory:" ascii
    condition:
        uint16(0) == 0x5a4d and filesize > 4MB and
        1 of ($x*) or 3 of ($s*)