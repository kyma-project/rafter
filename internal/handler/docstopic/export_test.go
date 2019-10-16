package docstopic

import "github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"

func FindSource() func(slice []v1beta1.Source, sourceName v1beta1.DocsTopicSourceName, sourceType v1beta1.DocsTopicSourceType) *v1beta1.Source {
	return findSource
}
