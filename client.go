// Package xrgrpc is a gRPC Client library for Cisco IOS XR devices. It
// exposes different RPC's to manage the device(s). The objective is
// to have a single interface to retrive info from the device, apply configs
// to it, generate telemetry streams and program the RIB/FIB.
package xrgrpc

import (
	"fmt"
	"io"
	"strconv"
	"time"

	pb "github.com/nleiva/xrgrpc/proto/ems"
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

// Devices identifies a list of gRPC targets
type Devices struct {
	Routers []CiscoGrpcClient
}

// NewDevices is a Devices constructor
func NewDevices() *Devices {
	return new(Devices)
}

// Provides the user/password for the connection. It implements
// the PerRPCCredentials interface.
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
func Connect(xr CiscoGrpcClient) (*grpc.ClientConn, context.Context, error) {
	// opts holds the config options to set up the connection.
	var opts []grpc.DialOption

	// creds provides the TLS credentials from the input certificate file.
	creds, err := credentials.NewClientTLSFromFile(xr.Creds, xr.Options)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to construct TLS credentialst")
	}
	// Add TLS credentials to config options array.
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// Add gRPC timeout to config options array.
	//opts = append(opts, grpc.WithTimeout(time.Second*time.Duration(xr.Timeout)))
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(xr.Timeout))

	// Add user/password to config options array.
	opts = append(opts, grpc.WithPerRPCCredentials(&loginCreds{
		Username: xr.User,
		Password: xr.Password}))

	// conn represents a client connection to an RPC server (target).
	conn, err := grpc.DialContext(ctx, xr.Host, opts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Fail to dial to target")
	}
	return conn, ctx, err
}

// ShowCmdTextOutput returns the output of a CLI show commands as text.
func ShowCmdTextOutput(ctx context.Context, conn *grpc.ClientConn, cli string, id int64) (string, error) {
	var s string
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
			si := strconv.FormatInt(id, 10)
			return s, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
		}
		if len(r.Output) > 0 {
			s += r.Output
		}
	}
}

// ShowCmdJSONOutput returns the output of a CLI show commands
// as a JSON structure output.
func ShowCmdJSONOutput(ctx context.Context, conn *grpc.ClientConn, cli string, id int64) (string, error) {
	var s string
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
			si := strconv.FormatInt(id, 10)
			return s, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
		}
		if len(r.Jsonoutput) > 0 {
			s += r.Jsonoutput
		}
	}
}

// GetConfig returns the config for a specif YANG path elments
// descibed in 'js'.
func GetConfig(ctx context.Context, conn *grpc.ClientConn, js string, id int64) (string, error) {
	var s string
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
			si := strconv.FormatInt(id, 10)
			return s, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
		}
		if len(r.Yangjson) > 0 {
			s += r.Yangjson
		}
	}
}

// CLIConfig configs the target with CLI commands descibed in 'cli'.
func CLIConfig(ctx context.Context, conn *grpc.ClientConn, cli string, id int64) error {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.CliConfigArgs{ReqId: id, Cli: cli}

	// 'r' is the result that comes back from the target.
	r, err := c.CliConfig(ctx, &a)
	if err != nil {
		return errors.Wrap(err, "gRPC CliConfig failed")
	}
	if len(r.Errors) != 0 {
		si := strconv.FormatInt(id, 10)
		return fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
	}
	return err
}

// CommitConfig commits a config. Need to clarify its use-case.
func CommitConfig(ctx context.Context, conn *grpc.ClientConn, id int64) (string, error) {
	var s string
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)
	si := strconv.FormatInt(id, 10)
	// Commit metadata
	m := pb.CommitMsg{Label: "gRPC id: " + si, Comment: "gRPC commit id: " + si}

	// 'a' is the object we send to the router via the stub.
	a := pb.CommitArgs{Msg: &m, ReqId: id}

	// 'r' is the result that comes back from the target.
	r, err := c.CommitConfig(context.Background(), &a)
	if err != nil {
		return s, errors.Wrap(err, "gRPC CommitConfig failed")
	}
	if len(r.Errors) != 0 {
		return s, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
	}
	// What about r.ResReqId. Seems to equal to id sent.
	return r.Result.String(), err
}

// DiscardConfig deletes configs with ID 'id' on the target.
// Need to clarify its use-case.
func DiscardConfig(ctx context.Context, conn *grpc.ClientConn, id int64) (int64, error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.DiscardChangesArgs{ReqId: id}

	// 'r' is the result that comes back from the target.
	r, err := c.ConfigDiscardChanges(context.Background(), &a)
	if err != nil {
		return -1, errors.Wrap(err, "gRPC ConfigDiscardChanges failed")
	}
	if len(r.Errors) != 0 {
		si := strconv.FormatInt(id, 10)
		return -1, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
	}
	return r.ResReqId, nil
}

// MergeConfig configs the target with YANG/JSON config specified in 'js'.
func MergeConfig(ctx context.Context, conn *grpc.ClientConn, js string, id int64) (int64, error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.ConfigArgs{ReqId: id, Yangjson: js}

	// 'r' is the result that comes back from the target.
	r, err := c.MergeConfig(context.Background(), &a)
	if err != nil {
		return -1, errors.Wrap(err, "gRPC MergeConfig failed")
	}
	if len(r.Errors) != 0 {
		si := strconv.FormatInt(id, 10)
		return -1, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
	}
	return r.ResReqId, nil
}

// DeleteConfig removes the config config specified in 'js'
// on the target device.
func DeleteConfig(ctx context.Context, conn *grpc.ClientConn, js string, id int64) (int64, error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.ConfigArgs{ReqId: id, Yangjson: js}

	// 'r' is the result that comes back from the target.
	r, err := c.DeleteConfig(context.Background(), &a)
	if err != nil {
		return -1, errors.Wrap(err, "gRPC DeleteConfig failed")
	}
	if len(r.Errors) != 0 {
		si := strconv.FormatInt(id, 10)
		return -1, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
	}
	return r.ResReqId, nil
}

// ReplaceConfig replaces the config specified in 'js' on
// the target device.
func ReplaceConfig(ctx context.Context, conn *grpc.ClientConn, js string, id int64) (int64, error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.ConfigArgs{ReqId: id, Yangjson: js}

	// 'r' is the result that comes back from the target.
	r, err := c.ReplaceConfig(ctx, &a)
	if err != nil {
		return -1, errors.Wrap(err, "gRPC ReplaceConfig failed")
	}
	if len(r.Errors) != 0 {
		si := strconv.FormatInt(id, 10)
		return -1, fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
	}
	return r.ResReqId, nil
}

// GetSubscription follows the Channel Generator Pattern, it
// returns a channel where the Streaming Telemetry data is received
func GetSubscription(ctx context.Context, conn *grpc.ClientConn, p string, id int64, e int64) (chan []byte, error) {
	// 'c' is the gRPC stub.
	c := pb.NewGRPCConfigOperClient(conn)
	// 'b' is the bytes channel where Telemetry is received
	b := make(chan []byte)
	// 'a' is the object we send to the router via the stub.
	a := pb.CreateSubsArgs{ReqId: id, Encode: e, Subidstr: p}

	// 'r' is the result that comes back from the target.
	st, err := c.CreateSubs(ctx, &a)
	if err != nil {
		return b, errors.Wrap(err, "gRPC CreateSubs failed")
	}

	// TODO: Review the logic.
	go func() {
		for {
			r, err := st.Recv()
			if err == io.EOF {
				close(b)
				break
			}
			if len(r.Errors) != 0 {
				si := strconv.FormatInt(id, 10)
				err = fmt.Errorf("Error triggered by remote host for ReqId: %s; %s", si, r.Errors)
				close(b)
				break
			}
			select {
			case <-ctx.Done():
				close(b)
				break
			case b <- r.GetData():
				continue
			}
		}
	}()
	return b, err
}
