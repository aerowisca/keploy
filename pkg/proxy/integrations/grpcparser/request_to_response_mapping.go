package grpcparser

import (
	"fmt"

	"go.keploy.io/server/pkg/models"
)

const (
	KLabelForAuthority = ":authority"
	KLabelForMethod    = ":method"
	KLabelForPath      = ":path"
	KLabelForScheme    = ":http"

	KLabelForContentType = "content-type"
)

func FilterMocksRelatedToGrpc(mocks []*models.Mock) []*models.Mock {
	var res []*models.Mock
	for _, mock := range mocks {
		if mock != nil && mock.Kind == models.GRPC_EXPORT && mock.Spec.GRPCReq != nil && mock.Spec.GRPCResp != nil {
			res = append(res, mock)
		}
	}
	return res
}

func FilterMocksBasedOnGrpcRequest(grpcReq models.GrpcReq, mocks []*models.Mock) *models.Mock {
	grpcMocks := FilterMocksRelatedToGrpc(mocks)
	for _, mock := range grpcMocks {
		have := mock.Spec.GRPCReq
		// Investigate pseudo headers.
		if have.Headers.PseudoHeaders[KLabelForAuthority] != grpcReq.Headers.PseudoHeaders[KLabelForAuthority] {
			continue
		}
		if have.Headers.PseudoHeaders[KLabelForMethod] != grpcReq.Headers.PseudoHeaders[KLabelForMethod] {
			continue
		}
		if have.Headers.PseudoHeaders[KLabelForPath] != grpcReq.Headers.PseudoHeaders[KLabelForPath] {
			continue
		}
		if have.Headers.PseudoHeaders[KLabelForScheme] != grpcReq.Headers.PseudoHeaders[KLabelForScheme] {
			continue
		}

		// Investigate ordinary headers.
		if have.Headers.OrdinaryHeaders[KLabelForContentType] != grpcReq.Headers.OrdinaryHeaders[KLabelForContentType] {
			continue
		}

		if len(have.Body) != len(grpcReq.Body) {
			fmt.Println("length not equal of the prefixedLength arrays ", len(have.Body), len(grpcReq.Body))
			continue
		} else {
			count := 0
			for i, body := range have.Body {
				// Investigate the compression flag.
				if body.CompressionFlag != grpcReq.Body[i].CompressionFlag {
					count = 0
					break
				}

				// Investigate the decoded data.
				if body.DecodedData != grpcReq.Body[i].DecodedData {
					count = 0
					break
				}
				count += 1
			}
			if count == len(have.Body) {
				return mock
			} else {
				continue
			}
		}

		//return mock
	}

	return nil
}
