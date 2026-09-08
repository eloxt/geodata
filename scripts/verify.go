// Validate every generated SRS and ensure all source categories were converted.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: verify INPUT_DIRECTORY OUTPUT_DIRECTORY")
	}
	expected := map[string]bool{}
	for _, kind := range []string{"geoip", "geosite"} {
		data, err := os.ReadFile(filepath.Join(os.Args[1], kind+".dat"))
		if err != nil {
			log.Fatal(err)
		}
		if kind == "geoip" {
			var list routercommon.GeoIPList
			if err = proto.Unmarshal(data, &list); err != nil {
				log.Fatal(err)
			}
			if len(list.Entry) == 0 {
				log.Fatal("empty geoip source")
			}
			for _, entry := range list.Entry {
				if entry.InverseMatch {
					log.Fatal("converter does not support inverse GeoIP entries")
				}
				expected["geoip/"+strings.ToLower(entry.CountryCode)+".srs"] = true
			}
		} else {
			var list routercommon.GeoSiteList
			if err = proto.Unmarshal(data, &list); err != nil {
				log.Fatal(err)
			}
			if len(list.Entry) == 0 {
				log.Fatal("empty geosite source")
			}
			for _, entry := range list.Entry {
				code := strings.ToLower(entry.CountryCode)
				expected["geosite/"+code+".srs"] = true
				for _, domain := range entry.Domain {
					for _, attribute := range domain.Attribute {
						expected["geosite/"+code+"@"+attribute.Key+".srs"] = true
					}
				}
			}
		}
	}
	for name := range expected {
		if strings.Contains(name, "..") || strings.Count(name, "/") != 1 || strings.Contains(name, "\\") {
			log.Fatalf("invalid source category: %q", name)
		}
		file, err := os.Open(filepath.Join(os.Args[2], name))
		if err != nil {
			log.Fatal(err)
		}
		rules, err := srs.Read(file, true)
		file.Close()
		if err != nil {
			log.Fatalf("%s: %v", name, err)
		}
		if len(rules.Rules) == 0 {
			log.Fatalf("%s: no rules", name)
		}
	}
	files, err := filepath.Glob(filepath.Join(os.Args[2], "*", "*.srs"))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) != len(expected) {
		log.Fatalf("category count mismatch: got %d, expected %d", len(files), len(expected))
	}
	fmt.Printf("Validated %d SRS files, including GeoSite attribute categories.\n", len(files))
}
