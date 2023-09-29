// Generates dns deploy scenario step, with 256 dns contracts, 1 per shard.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	mc "github.com/kalyan3104/dme-vm-util/test-util/mandos/controller"
	mj "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/model"
	mjparse "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/parse"
	mjwrite "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/write"
)

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../../test")
	return arwenTestRoot
}

type testGenerator struct {
	mandosParser      mjparse.Parser
	generatedScenario *mj.Scenario
}

func (tg *testGenerator) addStep(stepSnippet string) {
	step, err := tg.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		panic(err)
	}
	tg.generatedScenario.Steps = append(tg.generatedScenario.Steps, step)
}

func main() {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"dns.wasm",
			filepath.Join(getTestRoot(), "dns/dns.wasm"))
	tg := &testGenerator{
		mandosParser: mjparse.Parser{
			FileResolver: fileResolver,
		},
		generatedScenario: &mj.Scenario{
			Name: "dns test",
		},
	}

	newAddressesSnippets := ""
	for shard := 0; shard < 256; shard++ {
		if shard > 0 {
			newAddressesSnippets += ","
		}
		newAddressesSnippets += fmt.Sprintf(`{
				"creatorAddress": "''dns_owner_______________________",
				"creatorNonce": "0x%02x",
				"newAddress": "''dns____________________________|0x%02x"
			}`,
			shard,
			shard)
	}
	tg.addStep(fmt.Sprintf(`
			{
				"step": "setState",
				"accounts": {
					"''dns_owner_______________________": {
						"nonce": "0",
						"balance": "0",
						"storage": {},
						"code": ""
					}
				},
				"newAddresses": [
					%s
				]
			}`,
		newAddressesSnippets))

	for shard := 0; shard < 256; shard++ {
		tg.addStep(fmt.Sprintf(`
			{
				"step": "scDeploy",
				"txId": "deploy-0x%02x",
				"tx": {
					"from": "''dns_owner_______________________",
					"value": "0",
					"contractCode": "file:dns.wasm",
					"arguments": [ "123,000" ],
					"gasLimit": "100,000",
					"gasPrice": "0"
				},
				"expect": {
					"out": [],
					"status": "",
					"logs": [],
					"gas": "*",
					"refund": "*"
				}
			}`,
			shard))
	}

	// save
	serialized := mjwrite.ScenarioToJSONString(tg.generatedScenario)
	err := ioutil.WriteFile(
		filepath.Join(getTestRoot(), "dns/dns_init.steps.json"),
		[]byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
