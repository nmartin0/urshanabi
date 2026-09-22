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
