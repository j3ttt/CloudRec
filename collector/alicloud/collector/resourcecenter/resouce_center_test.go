package resourcecenter

import (
	"github.com/cloudrec/alicloud/collector"
	"github.com/core-sdk/constant"
	"github.com/core-sdk/log"
	"github.com/core-sdk/schema"
	"os"
	"testing"
)

var GetTestAccount = func() (res []schema.CloudAccount) {
	testAccount := schema.CloudAccount{
		CloudAccountId: "test-account",
		CommonCloudAccountAuthParam: schema.CommonCloudAccountAuthParam{
			AK: os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"),
			SK: os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"),
		},
	}

	res = append(res, testAccount)

	return res
}

func TestCloudCenterResource(t *testing.T) {
	p := schema.GetInstance(schema.PlatformConfig{
		Name: string(constant.AlibabaCloud),
		Resources: []schema.Resource{
			GeCloudCenterResource(),
		},

		Service: &collector.Services{},
		DefaultRegions: []string{
			"cn-shanghai",
			"ap-southeast-1",
		},
		DefaultCloudAccounts: GetTestAccount(),
	})

	if err := schema.RunExecutor(p); err != nil {
		log.GetWLogger().Error(err.Error())
	}
}
