package model

// Sys model
type Sys struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	LinkType    string `json:"linkType,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	UpdatedBy   *Sys   `json:"updatedBy,omitempty"`
	Version     int    `json:"version,omitempty"`
	Revision    int    `json:"revision,omitempty"`
	ContentType *struct {
		Sys Sys `json:"sys,omitempty"`
	} `json:"contentType,omitempty"`
	Space *struct {
		Sys Sys `json:"sys,omitempty"`
	} `json:"space,omitempty"`
	FirstPublishedAt string `json:"firstPublishedAt,omitempty"`
	PublishedCounter int    `json:"publishedCounter,omitempty"`
	PublishedAt      string `json:"publishedAt,omitempty"`
	PublishedBy      *Sys   `json:"publishedBy,omitempty"`
	PublishedVersion int    `json:"publishedVersion,omitempty"`
	ArchivedAt       string `json:"archivedAt,omitempty"`
	ArchivedBy       *Sys   `json:"archivedBy,omitempty"`
	ArchivedVersion  int    `json:"archivedVersion,omitempty"`
	AppDefinition    *struct {
		Sys Sys `json:"sys,omitempty"`
	} `json:"appDefinition,omitempty"`
}
