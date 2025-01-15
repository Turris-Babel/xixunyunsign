package cmd

type RequestPayload struct {
	Contents []Content `json:"contents"`
}
type Response struct {
	Candidates    []Candidate `json:"candidates"`
	UsageMetadata interface{} `json:"usageMetadata"` // 根据需要定义具体类型
	ModelVersion  string      `json:"modelVersion"`
}
type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason"`
	AvgLogprobs  float64 `json:"avgLogprobs"`
}

type Content struct {
	Parts []ContentPart `json:"parts"`
	Role  string        `json:"role"`
}
type ContentPart struct {
	Text string `json:"text"`
}

type UserInfo struct {
	account  string
	password string
	schoolID string
	token    string
}

type signconfig struct {
	address      string
	address_name string
	latitude     string
	longitude    string
	remark       string
	comment      string
	province     string
	city         string
	debug        bool
}
