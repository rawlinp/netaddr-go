package netaddr

import "testing"

func Test_ParseIPv6(t *testing.T) {
	cases := []struct {
		given     string
		hi64      uint64
		lo64      uint64
		expectErr bool
	}{
		{" :: ", 0, 0, false},
		{"::0", 0, 0, false},
		{"::1", 0, 1, false},
		{"fe80::", 0xfe80000000000000, 0, false},
		//{"::ffff:192.168.1.1",0,0xffffc0a80101,false}, // ipv4 mapped
		{"fe80::1::", 0, 0, true},
		{"::fe80::", 0, 0, true},
		{"0:0:0:0:0:0:0:0:1", 0, 0, true},
		{"::0:0:0:0:0:0:1", 0, 0, true},
		{"1:1:1:1:1:1:1::", 0, 0, true},
		{"1:1:1:1:1:1:1::", 0, 0, true},
		{"fec0", 0, 0, true},
		{"fec0:::1", 0, 0, true},
	}

	for _, c := range cases {
		ip, err := ParseIPv6(c.given)
		if err != nil {
			if !c.expectErr {
				t.Errorf("ParseIPv6(%s) unexpected parse error: %s", c.given, err.Error())
			}
			continue
		}

		if c.expectErr {
			t.Errorf("ParseIPv6(%s) expected error but none raised", c.given)
			continue
		}

		if ip.netId != c.hi64 || ip.hostId != c.lo64 {
			t.Errorf("ParseIPv6(%s)  Expect: %x%x  Result: %x%x", c.given, c.hi64, c.lo64, ip.netId, ip.hostId)
		}
	}
}

func Test_IPv6_Cmp(t *testing.T) {
	cases := []struct {
		ip1 string
		ip2 string
		res int
	}{
		{"::", "::1", -1},  // hostId numerically less
		{"::1", "::", 1},   // hostId numerically greater
		{"::1", "::1", 0},  // hostId eq
		{"1::", "2::", -1}, // netId numerically less
		{"2::", "1::", 1},  // netId numerically greater
		{"1::", "1::", 0},  // netId eq
	}

	for _, c := range cases {
		ip1, _ := ParseIPv6(c.ip1)
		ip2, _ := ParseIPv6(c.ip2)

		if res, _ := ip1.Cmp(ip2); res != c.res {
			t.Errorf("%s.Cmp(%s) Expect: %d  Result: %d", ip1, ip2, c.res, res)
		}
	}
}

func Test_IPv6_Long(t *testing.T) {
	cases := []struct {
		given  string
		expect string
	}{
		{"::", "0000:0000:0000:0000:0000:0000:0000:0000"},
		{"1::", "0001:0000:0000:0000:0000:0000:0000:0000"},
		{"1000::", "1000:0000:0000:0000:0000:0000:0000:0000"},
	}

	for _, c := range cases {
		ip, _ := ParseIPv6(c.given)
		long := ip.Long()
		if long != c.expect {
			t.Errorf("%s.Long() Expect: %s  Result: %s", c.given, c.expect, long)
		}
	}
}

func Test_IPv6_String(t *testing.T) {
	cases := []struct {
		given  string
		expect string
	}{
		{"0:0:0:0:0:0:0:0", "::"},
		{"1:0:0:0:0:0:0:0", "1::"},
		{"0:1:0:0:0:0:0:0", "0:1::"},
		{"0:0:1:0:0:0:0:0", "0:0:1::"},
		{"0:0:0:1:0:0:0:0", "0:0:0:1::"},
		{"0:0:0:0:1:0:0:0", "::1:0:0:0"},
		{"0:0:0:0:0:1:0:0", "::1:0:0"},
		{"0:0:0:0:0:0:1:0", "::1:0"},
		{"0:0:0:0:0:0:0:1", "::1"},
		
		{"1:0:0:0:0:0:0:1", "1::1"},
		{"1:1:0:0:0:0:0:1", "1:1::1"},
		{"1:0:1:0:0:0:0:1", "1:0:1::1"},
		{"1:0:0:1:0:0:0:1", "1:0:0:1::1"},
		{"1:0:0:0:1:0:0:1", "1::1:0:0:1"},
		{"1:0:0:0:0:1:0:1", "1::1:0:1"},
		{"1:0:0:0:0:0:1:1", "1::1:1"},
	}

	for _, c := range cases {
		ip, _ := ParseIPv6(c.given)
		short := ip.String()
		if short != c.expect {
			t.Errorf("%s.String() Expect: %s  Result: %s", c.given, c.expect, short)
		}
	}
}
