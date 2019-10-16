package route

import "github.com/kyma-project/rafter/pkg/extractor"

func (h *ExtractHandler) SetMetadataExtractor(metadataExtractor extractor.Extractor) {
	h.metadataExtractor = metadataExtractor
}
