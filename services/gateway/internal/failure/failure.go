// Package failure describes errors from other services for the log
// without revealing what they must not (roadmap R-06). Every
// application error carries an ErrorDetail; its reason and safe
// arguments are logged, and its unsafe arguments -- anything that could
// hold customer data, credentials or query text -- never are. Nor is the
// producer's free-text message, which could echo any of them. An error
// with no detail comes from the transport itself, such as a refused
// connection, so its message is logged.
package failure

import (
	"google.golang.org/grpc/codes"  // name-ok
	"google.golang.org/grpc/status" // name-ok

	commonv1 "urshanabi/services/gateway/internal/gen/urshanabi/common/v1"
)

// Attrs returns what may be logged about err, as key-value pairs.
func Attrs(err error) []any {
	st, ok := status.FromError(err)
	if !ok {
		return []any{"error", err.Error()}
	}
	attrs := []any{"code", st.Code().String()}
	for _, d := range st.Details() {
		detail, ok := d.(*commonv1.ErrorDetail)
		if !ok {
			continue
		}
		attrs = append(attrs, "reason", detail.GetReason())
		for _, a := range detail.GetSafeArgs() {
			attrs = append(attrs, "arg."+a.GetName(), a.GetValue())
		}
		return attrs
	}
	return append(attrs, "error", st.Message())
}

// The four kinds a caller may be told about (roadmap R-22), in our own
// vocabulary: nothing else ever reaches a caller, so no answer reveals
// the language a service is written in or the library it failed inside.
const (
	Unreachable   = "unreachable"
	Refused       = "refused"
	Empty         = "empty"
	Misconfigured = "misconfigured"
)

// Kind says which of the four kinds err is. A producer that states its
// own kind is believed; otherwise the transport's code decides.
func Kind(err error) string {
	st, ok := status.FromError(err)
	if !ok {
		return Misconfigured
	}
	for _, d := range st.Details() {
		if detail, ok := d.(*commonv1.ErrorDetail); ok {
			switch detail.GetKind() {
			case commonv1.FailureKind_FAILURE_KIND_UNREACHABLE:
				return Unreachable
			case commonv1.FailureKind_FAILURE_KIND_REFUSED:
				return Refused
			case commonv1.FailureKind_FAILURE_KIND_EMPTY:
				return Empty
			case commonv1.FailureKind_FAILURE_KIND_MISCONFIGURED:
				return Misconfigured
			case commonv1.FailureKind_FAILURE_KIND_UNSPECIFIED:
			}
		}
	}
	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.Canceled, codes.Aborted:
		return Unreachable
	case codes.PermissionDenied, codes.Unauthenticated, codes.InvalidArgument,
		codes.FailedPrecondition, codes.OutOfRange, codes.ResourceExhausted, codes.AlreadyExists:
		return Refused
	case codes.NotFound:
		return Empty
	default:
		return Misconfigured
	}
}
