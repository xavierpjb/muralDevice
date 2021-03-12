package artifact

import "time"

// ArtifactModel represent the model to imitate an artifact (file/mp4 etc)
type ArtifactModel struct {
	Username       string
	File           string
	Type           string
	UploadDateTime time.Time
	Caption        *string
}

// IsPersistable checks that the properties need to persist an artifact are all present
func (a ArtifactModel) IsPersistable() bool {
	return a.File != "" && a.Type != "" && a.Username != "" && !a.UploadDateTime.IsZero()
}

// DeleteModel specifies the information needed to delete an artifact
type DeleteModel struct {
	Username string
	URL      string
}

// IsDeleteable checks for the params needed to make a deleted request
func (a DeleteModel) IsDeleteable() bool {
	return a.Username != "" && a.URL != ""
}
