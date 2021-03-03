package artifact

// ArtifactModel represent the model to imitate an artifact (file/mp4 etc)
type ArtifactModel struct {
	Username string
	File     string
	Type     string
}

// IsPersistable checks that the properties need to persist an artifact are all present
func (a ArtifactModel) IsPersistable() bool {
	return a.File != "" && a.Type != "" && a.Username != ""
}
