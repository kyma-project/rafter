package uploader

func (u *Uploader) PopulateErrors(errorsCh chan *UploadError) []UploadError {
	return u.populateErrors(errorsCh)
}

func (u *Uploader) PopulateResults(resultsCh chan *UploadResult) []UploadResult {
	return u.populateResults(resultsCh)
}

func (u *Uploader) NormalizeObjectName(dir, fileName string) string {
	return u.normalizeObjectName(dir, fileName)
}
