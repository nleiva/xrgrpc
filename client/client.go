/*
gRPC Client library
*/

package xrgrpc

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/nleiva/xrgrpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// CiscoGrpcClient identifies the parameters for gRPC session setup.
type CiscoGrpcClient struct {
	User     string
	Password string
	Host     string
	Creds    string
	Options  string
	Timeout  int
}

// NewCiscoGrpcClient is a CiscoGrpcClient constructor
func NewCiscoGrpcClient() *CiscoGrpcClient {
	return new(CiscoGrpcClient)
}

// Provides the user/password for the connection.
// It implements the PerRPCCredentials interface.
// https://godoc.org/google.golang.org/grpc/credentials#PerRPCCredentials.
type loginCreds struct {
	Username, Password string
}

// Method of the PerRPCCredentials interface.
func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": c.Username,
		"password": c.Password,
	}, nil
}

// Method of the PerRPCCredentials interface.
func (c *loginCreds) RequireTransportSecurity() bool {
	return true
}

// Connect will return a grpc.ClienConn to the target.
func Connect(xr CiscoGrpcClient) (conn *grpc.ClientConn, err error) {
	// opts holds the config options to set up the connection.
	var opts []grpc.DialOption

	// creds provides the TLS credentials from the input certificate file.
	creds, err := credentials.NewClientTLSFromFile(xr.Creds, xr.Options)
	if err != nil {
		log.Fatalf("Failed to construct TLS credentials: %v", err)
		return nil, err
	}
	// Add TLS credentials to config options array.
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// Add gRPC timeout to config options array.
	opts = append(opts, grpc.WithTimeout(time.Second*time.Duration(xr.Timeout)))

	// Add user/password to config options array.
	opts = append(opts, grpc.WithPerRPCCredentials(&loginCreds{
		Username: xr.User,
		Password: xr.Password}))

	// conn represents a client connection to an RPC server (target).
	conn, err = grpc.Dial(xr.Host, opts...)
	if err != nil {
		log.Fatalf("Fail to dial to target: %v", err)
		return nil, err
	}
	return conn, err
}

// ShowCmdTextOutput returns the output of a CLI show commands as text.
func ShowCmdTextOutput(conn *grpc.ClientConn, cli string, id int64) (s string, err error) {
	// 'client' is the gRPC stub.
	client := pb.NewGRPCExecClient(conn)

	// 'cliArgs' is the object we send to the router via the stub.
	cliArgs := pb.ShowCmdArgs{ReqId: id, Cli: cli}

	// 'stream' is the streamed result that comes back from the target.
	stream, err := client.ShowCmdTextOutput(context.Background(), &cliArgs)
	if err != nil {
		return s, err
	}

	for {
		// Loop through the responses in the stream until there is nothing left.
		reply, err := stream.Recv()
		if err == io.EOF {
			return s, nil
		}
		if len(reply.Errors) != 0 {
			err := errors.New("Error triggered by remote host")
			fmt.Printf(
				"ShowCmd: ReqId %d, received error: %s\n",
				id,
				reply.Errors)
			return s, err
		}
		if len(reply.Output) > 0 {
			s += reply.Output
		}
	}
}

// ShowCmdJSONOutput returns the output of a CLI show commands as a JSON structure output.
// A lot of code duplication (from ShowCmdTextOutput). Will improve this.
func ShowCmdJSONOutput(conn *grpc.ClientConn, cli string, id int64) (s string, err error) {
	// 'client' is the gRPC stub.
	client := pb.NewGRPCExecClient(conn)

	// 'cliArgs' is the object we send to the router via the stub.
	cliArgs := pb.ShowCmdArgs{ReqId: id, Cli: cli}

	// 'stream' is the streamed result that comes back from the target.
	stream, err := client.ShowCmdJSONOutput(context.Background(), &cliArgs)
	if err != nil {
		return s, err
	}

	for {
		// Loop through the responses in the stream until there is nothing left.
		reply, err := stream.Recv()
		if err == io.EOF {
			return s, nil
		}
		if len(reply.Errors) != 0 {
			err := errors.New("Error triggered by remote host")
			fmt.Printf(
				"ShowCmd: ReqId %d, received error: %s\n",
				id,
				reply.Errors)
			return s, err
		}
		if len(reply.Jsonoutput) > 0 {
			s += reply.Jsonoutput
		}
	}
}

// GetConfig returns the config for a specif YANG path elments descibed in js
// A lot of code duplication (from ShowCmdTextOutput). Will improve this.
func GetConfig(conn *grpc.ClientConn, js string, id int64) (s string, err error) {
	// 'client' is the gRPC stub.
	client := pb.NewGRPCConfigOperClient(conn)

	// 'jsonArgs' is the object we send to the router via the stub.
	jsonArgs := pb.ConfigGetArgs{ReqId: id, Yangpathjson: js}

	// 'stream' is the streamed result that comes back from the target.
	stream, err := client.GetConfig(context.Background(), &jsonArgs)
	if err != nil {
		return s, err
	}

	for {
		// Loop through the responses in the stream until there is nothing left.
		reply, err := stream.Recv()
		if err == io.EOF {
			return s, nil
		}
		if len(reply.Errors) != 0 {
			err := errors.New("Error triggered by remote host")
			fmt.Printf(
				"GetConfig: ReqId %d, received error: %s\n",
				id,
				reply.Errors)
			return s, err
		}
		if len(reply.Yangjson) > 0 {
			s += reply.Yangjson
		}
	}
}
