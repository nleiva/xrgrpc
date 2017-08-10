// Big TODO
package xrgrpc

import (
	"net"
	"strings"
	"testing"
	"time"

	pb "github.com/nleiva/xrgrpc/proto/ems"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	defaultAddr = "localhost"
	defaultPort = ":57344"
	defaultUser = "test"
	defaultPass = "test"
	defaultCert = "test/cert.pem"
	defaultKey  = "test/key.pem"
)

// server implements the GRPCExecServer interface
type server struct{}

func (s *server) ShowCmdTextOutput(a *pb.ShowCmdArgs, g pb.GRPCExec_ShowCmdTextOutputServer) error {
	return nil
}

func (s *server) ShowCmdJSONOutput(a *pb.ShowCmdArgs, g pb.GRPCExec_ShowCmdJSONOutputServer) error {
	return nil
}

//type GRPCExec_ShowCmdTextOutputServer interface {
//	Send(*ShowCmdTextReply) error
//	grpc.ServerStream
//}

//type ShowCmdTextReply struct {
//	ResReqId int64  `protobuf:"varint,1,opt,name=ResReqId" json:"ResReqId,omitempty"`
//	Output   string `protobuf:"bytes,2,opt,name=output" json:"output,omitempty"`
//	Errors   string `protobuf:"bytes,3,opt,name=errors" json:"errors,omitempty"`
//}

// streamInterceptor to authenticate incoming gRPC stream connections
func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authorize(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

// unaryInterceptor to authenticate incoming gRPC unary connections
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorize(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// Validates username and password
func authorize(ctx context.Context) error {
	if md, ok := metadata.FromContext(ctx); ok {
		if len(md["username"]) > 0 && md["username"][0] == defaultUser &&
			len(md["password"]) > 0 && md["password"][0] == defaultPass {
			return nil
		}

		return errors.New("incorrect username/password")
	}

	return errors.New("empty metadata")
}

func Server(t *testing.T) *grpc.Server {
	lis, err := net.Listen("tcp", defaultPort)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile(defaultCert, defaultKey)
	if err != nil {
		t.Fatalf("failed to construct TLS credentialst: %v", err)
	}
	// var opts []grpc.ServerOption
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	// pb.RegisterGRPCExecServer(s, &server{})
	go func() {
		err := s.Serve(lis)
		// Serve always returns a non-nil error :-(
		if strings.Contains(err.Error(), "use of closed network connection") {
			return
		}
		if err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()
	return s
}

// We are mainly validating we can connect to a gRPC Server
func TestConnect(t *testing.T) {
	// _ = context.Background()
	x := CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  5,
	}
	s := Server(t)

	conn, ctx, err := Connect(x)
	// No Error if it keeps re-trying.
	if err != nil {
		t.Fatalf("could not setup a client connection to %v", x.Host)
	}
	select {
	case <-ctx.Done():
		t.Fatalf("could not setup a client connection to %v in under 2 seconds", x.Host)
	// Just wait for 3 seconds for this local connection to be setup.
	case <-time.After(3 * time.Second):
		break
	}
	// Connection won't fail until it timeouts. It re-attempt to connect until this happens.
	// Can initially timeout because of the WithTimeout option hard-coded to two seconds
	// or after an overal timeout of 'x.Timeout'
	conn.Close()
	s.GracefulStop()
}
