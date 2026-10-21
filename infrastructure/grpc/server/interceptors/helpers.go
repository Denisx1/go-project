package interceptors

import "google.golang.org/grpc/codes"

func isServerFault(c codes.Code) bool {
	switch c {
	case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
		return true
	}
	return false
}
