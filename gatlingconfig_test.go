package libgatlingconfig

import (
	"fmt"
	"os"
	"testing"
)

func TestGetSingleGatlingConfig(t *testing.T) {
	err := os.Setenv("CONFIGSERVERURL", "https://appconfigwriter.wjapptokenconfig.workers.dev")
	if err != nil {
		t.Errorf("err %v", err)
		return
	}
	err = os.Setenv("CONFIGSERVERTOKEN", "")
	if err != nil {
		t.Errorf("err %v", err)
		return
	}
	err = os.Setenv("APPTOKEN", "LibGatlingConfigTest-001")
	if err != nil {
		t.Errorf("err %v", err)
		return
	}

	gConfig := GetSingleGatlingConfig()
	err = gConfig.Initialize("GatlingConfigTest")
	if err != nil {
		t.Errorf("gconfig init err %v", err)
		return
	}

	data1 := gConfig.Get("config1")

	if data1 != "config2" {
		t.Errorf("value err should be config2")
		return
	}

	fmt.Println("good: ", data1)
}
