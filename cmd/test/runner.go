package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	am "github.com/kalyan3104/dme-vm-go/arwenmandos"
	mc "github.com/kalyan3104/dme-vm-util/test-util/mandos/controller"
)

func resolveArgument(arg string) (string, bool, error) {
	fi, err := os.Stat(arg)
	if os.IsNotExist(err) {
		exeDir, err := os.Getwd()
		if err != nil {
			return "", false, err
		}
		arg = filepath.Join(exeDir, arg)
		fmt.Println(arg)
		fi, err = os.Stat(arg)
	}
	if err != nil {
		return "", false, err
	}
	return arg, fi.IsDir(), nil
}

func main() {
	if len(os.Args) != 2 {
		panic("One argument expected - the path to the json test.")
	}

	jsonFilePath, isDir, err := resolveArgument(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// init
	executor, err := am.NewArwenTestExecutor()
	if err != nil {
		panic("Could not instantiate Arwen VM")
	}

	// execute
	switch {
	case isDir:
		runner := mc.NewScenarioRunner(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunAllJSONScenariosInDirectory(
			jsonFilePath,
			"",
			".scen.json",
			[]string{})
	case strings.HasSuffix(jsonFilePath, ".scen.json"):
		runner := mc.NewScenarioRunner(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONScenario(jsonFilePath)
	default:
		runner := mc.NewTestRunner(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONTest(jsonFilePath)
	}

	// print result
	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
