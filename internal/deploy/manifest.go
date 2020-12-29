package deploy

// Manifest represents a hermes stack manifest consisting of: hermes init version, hermes deploy version, associated resources and deployed files
type Manifest struct {
	InitVersion   string        `json:"init_version"`
	DeployVersion string        `json:"deploy_version"`
	Files         []File        `json:"files"`
	CloudFront    string        `json:"cloudfront"`
	Bucket        string        `json:"bucket"`
	Domain        string        `json:"domain"`
	DomainAliases []string      `json:"domain_aliases"`
	EdgeHandlers  []EdgeHandler `json:"edge_handlers"`
}

// File represents a file in a bucket
type File struct {
	Key      string `json:"key"`
	Checksum string `json:"checksum"` // Algorithm: SHA256
	Size     int64  `json:"size"`
}

// EdgeHandler represents a lambda@edge definition
type EdgeHandler struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	Type   string `json:"type"` // [ORIGIN_REQUEST, ORIGIN_RESPONSE, VIEWER_REQUEST, VIEWER_RESPONSE]
}
