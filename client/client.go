/*
gRPC Client library
*/

package xrgrpc

import (
	"io"
	"strconv"
	"time"

	pb "github.com/nleiva/xrgrpc/proto"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "Failed to construct TLS credentialst")
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
		return nil, errors.Wrap(err, "Fail to dial to target")
	}
	return conn, err
}

// ShowCmdTextOutput returns the output of a CLI show commands as text.
func ShowCmdTextOutput(conn *grpc.ClientConn, cli string, id int64) (s string, err error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCExecClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.ShowCmdArgs{ReqId: id, Cli: cli}

	// 'st' is the streamed result that comes back from the target.
	st, err := c.ShowCmdTextOutput(context.Background(), &a)
	if err != nil {
		return s, errors.Wrap(err, "gRPC ShowCmdTextOutput failed")
	}

	for {
		// Loop through the responses in the stream until there is nothing left.
		r, err := st.Recv()
		if err == io.EOF {
			return s, nil
		}
		if len(r.Errors) != 0 {
			return s, errors.New("Error triggered by remote host for ReqId: " + strconv.FormatInt(id, 10) + ": " + r.Errors)
		}
		if len(r.Output) > 0 {
			s += r.Output
		}
	}
}

// ShowCmdJSONOutput returns the output of a CLI show commands as a JSON structure output.
// A lot of code duplication (from ShowCmdTextOutput). Will improve this.
func ShowCmdJSONOutput(conn *grpc.ClientConn, cli string, id int64) (s string, err error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCExecClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.ShowCmdArgs{ReqId: id, Cli: cli}

	// 'st' is the streamed result that comes back from the target.
	st, err := c.ShowCmdJSONOutput(context.Background(), &a)
	if err != nil {
		return s, errors.Wrap(err, "gRPC ShowCmdJSONOutput failed")
	}

	for {
		// Loop through the responses in the stream until there is nothing left.
		r, err := st.Recv()
		if err == io.EOF {
			return s, nil
		}
		if len(r.Errors) != 0 {
			return s, errors.New("Error triggered by remote host for ReqId: " + strconv.FormatInt(id, 10) + ": " + r.Errors)
		}
		if len(r.Jsonoutput) > 0 {
			s += r.Jsonoutput
		}
	}
}

// GetConfig returns the config for a specif YANG path elments descibed in 'js'.
// A lot of code duplication (from ShowCmdTextOutput). Will improve this.
func GetConfig(conn *grpc.ClientConn, js string, id int64) (s string, err error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.ConfigGetArgs{ReqId: id, Yangpathjson: js}

	// 'st' is the streamed result that comes back from the target.
	st, err := c.GetConfig(context.Background(), &a)
	if err != nil {
		return s, errors.Wrap(err, "gRPC GetConfig failed")
	}

	for {
		// Loop through the responses in the stream until there is nothing left.
		r, err := st.Recv()
		if err == io.EOF {
			return s, nil
		}
		if len(r.Errors) != 0 {
			return s, errors.New("Error triggered by remote host for ReqId: " + strconv.FormatInt(id, 10) + ": " + r.Errors)
		}
		if len(r.Yangjson) > 0 {
			s += r.Yangjson
		}
	}
}

// CLIConfig configs the target with CLI commands descibed in 'cli'.
func CLIConfig(conn *grpc.ClientConn, cli string, id int64) error {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.CliConfigArgs{ReqId: id, Cli: cli}

	// 'r' is the result that comes back from the target.
	r, err := c.CliConfig(context.Background(), &a)
	if err != nil {
		return errors.Wrap(err, "gRPC CliConfig failed")
	}
	if len(r.Errors) != 0 {
		return errors.New("Error triggered by remote host for ReqId: " + strconv.FormatInt(id, 10) + ": " + r.Errors)
	}
	return nil
}

// CommitConfig commits the config submitted with id 'id'.
func CommitConfig(conn *grpc.ClientConn, id int64) (s string, err error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	m := pb.CommitMsg{Label: "", Comment: ""}

	// 'a' is the object we send to the router via the stub.
	a := pb.CommitArgs{Msg: &m, ReqId: id}

	// 'r' is the result that comes back from the target.
	r, err := c.CommitConfig(context.Background(), &a)
	if err != nil {
		return s, errors.Wrap(err, "gRPC CommitConfig failed")
	}
	if len(r.Errors) != 0 {
		return s, errors.New("Error triggered by remote host for ReqId: " + strconv.FormatInt(id, 10) + ": " + r.Errors)
	}
	// What about r.ResReqId?. Is it equal to ReqId?
	return r.Result.String(), nil
}
