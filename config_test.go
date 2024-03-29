package blocked

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/swoiow/dns_utils/loader"
	"github.com/swoiow/dns_utils/parsers"
)

const (
	rulesetPath = ".github/ruleset.txt"
	rulesetData = "rules.dat"
)

var defaultConfigs = NewConfigs()

func TestCreateCache(t *testing.T) {
	t.Skip()

	rules := &[]string{
		rulesetPath,
	}

	filter := bloom.NewWithEstimates(uint(defaultConfigs.Size), defaultConfigs.Rate)

	for _, rule := range *rules {
		m := loader.DetectMethods(rule)
		lines, err := m.LoadRules(m.StrictMode)
		if err != nil {
			panic(err)
		}
		addLines2filter(lines, filter)
	}

	wFile, err := os.Create(rulesetData)
	if err != nil {
		panic(err)
	}
	defer wFile.Close()

	num, err := filter.WriteTo(wFile)
	fmt.Printf("Saved %v about %v rules from filter.", num, filter.K())
}

func TestCacheByLocal(t *testing.T) {
	filter := bloom.NewWithEstimates(uint(defaultConfigs.Size), defaultConfigs.Rate)

	m := loader.DetectMethods(rulesetData)
	err := m.LoadCache(filter)
	if err != nil {
		panic(err)
	}

	type rule struct {
		name   string
		result bool
	}

	var items = []rule{
		{name: "baidu.com", result: false},
		{name: "reddit.com", result: false},
		{name: "xhscdn.com", result: false},
		{name: "*.xhscdn.com", result: false},
	}

	for _, tt := range items {
		t.Run(tt.name, func(t *testing.T) {
			if resp := filter.TestString(tt.name); resp != tt.result {
				t.Errorf("TestCache failed %v", tt.name)
			}
		})
	}
}

func TestCacheByFile(t *testing.T) {
	filter := bloom.NewWithEstimates(uint(defaultConfigs.Size), defaultConfigs.Rate)

	m := loader.DetectMethods(rulesetData)
	err := m.LoadCache(filter)
	if err != nil {
		panic(err)
	}

	type rule struct {
		name   string
		result bool
	}

	var items = []rule{}
	m2 := loader.DetectMethods(rulesetPath)

	lines, _ := m2.LoadRules(false)

	for _, line := range lines {
		result, domains := parsers.Parse(line, parsers.DomainParser)
		items = append(items, rule{name: strings.ToLower(domains[0]), result: result})
	}

	for _, tt := range items {
		t.Run(tt.name, func(t *testing.T) {
			if resp := filter.TestString(tt.name); resp != tt.result {
				t.Errorf("TestCache failed %v", tt.name)
			}
		})
	}
}
