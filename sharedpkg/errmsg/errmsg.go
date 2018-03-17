package errmsg

const (
    ErrMsgJsonResourceNotFound      string = "{\"error\":\"resource not found\"}"
    ErrMsgJsonResourceForbidden     string = "{\"error\":\"forbidden resource\"}"
    ErrMsgJsonInternalServiceIssue  string = "{\"error\":\"internal service error\"}"
    ErrMsgJsonUnallowedCountry      string = "{\"error\":\"uncovered region\"}"
    ErrMsgJsonInvalidInvitation     string = "{\"error\":\"invalid invitation code\"}"
    ErrMsgJsonUnsubmittedDevice     string = "{\"error\":\"invitation code used\"}"
)