package source

import "github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"

// returns index to the first source object from given slice with given source type
// or -1 if not found
func IndexByType(slice []v1beta1.Source, sourceType v1beta1.AssetGroupSourceType) int {
	for i, source := range slice {
		if source.Type != sourceType {
			continue
		}
		return i
	}
	return -1
}

// returns a copy of given slice that will not contain sources with given source type
func FilterByType(sources []v1beta1.Source, sourceType v1beta1.AssetGroupSourceType) []v1beta1.Source {
	var result []v1beta1.Source
	for _, source := range sources {
		if source.Type == sourceType {
			continue
		}
		result = append(result, source)
	}
	return result
}
