package failure

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"  // name-ok
	"google.golang.org/grpc/status" // name-ok

	commonv1 "urshanabi/services/gateway/internal/gen/urshanabi/common/v1"
)

func TestEachTransportCodeBecomesOneOfOurKinds(t *testing.T) {
	for code, want := range map[codes.Code]string{
		codes.Unavailable:       Unreachable,
		codes.DeadlineExceeded:  Unreachable,
		codes.Canceled:          Unreachable,
		codes.PermissionDenied:  Refused,
		codes.Unauthenticated:   Refused,
		codes.InvalidArgument:   Refused,
		codes.ResourceExhausted: Refused,
		codes.NotFound:          Empty,
		codes.Unimplemented:     Misconfigured,
		codes.Internal:          Misconfigured,
	} {
		if got := Kind(status.Error(code, "whatever the producer said")); got != want {
			t.Errorf("%v became %q, want %q", code, got, want)
		}
	}
}

func TestAProducerStatingItsOwnKindIsBelieved(t *testing.T) {
	st, err := status.New(codes.Internal, "whatever").WithDetails(&commonv1.ErrorDetail{
		Reason: "SOURCE_UNREACHABLE",
		Kind:   commonv1.FailureKind_FAILURE_KIND_UNREACHABLE,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := Kind(st.Err()); got != Unreachable {
		t.Errorf("got %q, want %q", got, Unreachable)
	}
}

func TestAFailureFromElsewhereNeedsAnOperator(t *testing.T) {
	if got := Kind(errors.New("something local")); got != Misconfigured {
		t.Errorf("got %q, want %q", got, Misconfigured)
	}
}
