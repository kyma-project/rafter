package testsuite

import (
	"fmt"
	"strings"

	"github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"
	"github.com/kyma-project/rafter/tests/asset-store/pkg/upload"
)

type assetData struct {
	Name string
	URL  string
	Mode v1beta1.AssetMode
}

func convertToAssetResourceDetails(response *upload.Response, prefix string) []assetData {
	var assets []assetData
	for _, file := range response.UploadedFiles {
		var mode v1beta1.AssetMode
		if strings.HasSuffix(file.FileName, ".tar.gz") || strings.HasSuffix(file.FileName, ".zip") {
			mode = v1beta1.AssetPackage
		} else {
			mode = v1beta1.AssetSingle
		}

		asset := assetData{
			Name: fmt.Sprintf("%s-%s", prefix, file.FileName),
			URL:  file.RemotePath,
			Mode: mode,
		}
		assets = append(assets, asset)
	}

	return assets
}
