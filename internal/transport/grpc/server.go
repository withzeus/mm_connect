package grpc

import (
	"context"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
	clientv1 "github.com/withzeus/id_contracts/gen/go/service/v1"
	"github.com/withzeus/mm_connect/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	clientv1.UnimplementedClientAuthServiceServer

	auth *service.AuthService
}

func New(auth *service.AuthService) *Server {
	return &Server{auth: auth}
}

func (s *Server) RegisterClient(ctx context.Context, req *clientv1.RegisterClientRequest) (*clientv1.RegisterClientResponse, error) {
	clientName := strings.TrimSpace(req.GetClientName())
	if clientName == "" {
		return nil, status.Error(codes.InvalidArgument, "Client name is required")
	}

	domain := strings.TrimSpace(strings.ToLower(req.GetDomain()))
	if domain == "" {
		return nil, status.Error(codes.InvalidArgument, "Domain name is required")
	}

	if !isValidDomainOrIP(domain) {
		return nil, status.Error(codes.InvalidArgument, "Domain must be a valid website host or IP address")
	}

	created, secret, err := s.auth.RegisterClient(ctx, clientName, domain)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register client: %v", err)
	}

	return &clientv1.RegisterClientResponse{
		ClientId:     created.UUID,
		ClientSecret: secret,
	}, nil

}

func isValidDomainOrIP(input string) bool {
	if ip := net.ParseIP(input); ip != nil {
		return true
	}

	testURL := input
	if !strings.HasPrefix(testURL, "http://") && !strings.HasPrefix(testURL, "https://") {
		testURL = "https://" + testURL
	}

	parsed, err := url.Parse(testURL)
	if err != nil {
		return false
	}

	host := parsed.Hostname()
	return host != "" && strings.Contains(host, ".")
}

func (s *Server) IssueToken(ctx context.Context, req *clientv1.IssueTokenRequest) (*clientv1.IssueTokenResponse, error) {
	clientIDStr := strings.TrimSpace(req.GetClientId())
	clientSecretStr := strings.TrimSpace(req.GetClientSecret())

	if clientIDStr == "" {
		return nil, status.Error(codes.InvalidArgument, "Client ID is required")
	}

	if _, err := uuid.Parse(clientIDStr); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Client ID")
	}

	if clientSecretStr == "" {
		return nil, status.Error(codes.InvalidArgument, "Client Secret is required")
	}

	if len(clientSecretStr) < 32 {
		return nil, status.Error(codes.InvalidArgument, "Invalid Client Secret")
	}

	token, err := s.auth.IssueToken(
		ctx,
		req.ClientId,
		req.ClientSecret,
	)

	if err != nil {
		return nil, err
	}

	return &clientv1.IssueTokenResponse{
		AccessToken: token,
	}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *clientv1.ValidateTokenRequest) (*clientv1.ValidateTokenResponse, error) {
	claims, err := s.auth.ValidateToken(
		ctx,
		req.Token,
	)

	if err != nil {
		return &clientv1.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &clientv1.ValidateTokenResponse{
		Valid:      true,
		ClientName: claims.ClientName,
		Scopes:     claims.Scopes,
	}, nil
}
