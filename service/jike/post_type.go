package jike

const (
	OriginalPost    SimplePostType = "originalPost"
	OfficialMessage                = "officialMessage"
)

const (
	RawOriginalPost    RawPostType = "ORIGINAL_POST"
	RawOfficialMessage             = "OFFICIAL_MESSAGE"
)

type (
	SimplePostType string

	RawPostType string
)

func (t RawPostType) ToSimpleType() SimplePostType {
	switch t {
	case RawOriginalPost:
		return OriginalPost
	case RawOfficialMessage:
		return OfficialMessage
	default:
		return OriginalPost
	}
}

func (t SimplePostType) ToRawType() RawPostType {
	switch t {
	case OriginalPost:
		return RawOriginalPost
	case OfficialMessage:
		return RawOfficialMessage
	default:
		return RawOriginalPost
	}
}

func (t RawPostType) String() string {
	return string(t)
}

func (t SimplePostType) String() string {
	return string(t)
}
