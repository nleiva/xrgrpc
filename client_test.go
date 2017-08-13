// Big TODO: current coverage: 26.9% of statements
package xrgrpc_test

import (
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"

	xr "github.com/nleiva/xrgrpc"
	pb "github.com/nleiva/xrgrpc/proto/ems"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	defaultAddr    = "localhost"
	defaultPort    = ":57344"
	defaultUser    = "test"
	defaultPass    = "test"
	defaultCert    = "test/cert.pem"
	defaultKey     = "test/key.pem"
	defaultCmd     = "show test"
	defaultYang    = "{\"Cisco-IOS-XR-test:tree\": [null]}"
	wrongCmd       = "show me the money"
	wrongCmdErr    = "wrong command"
	wrongYangErr   = "wrong YANG path"
	wrongCreds     = "incorrect username/password"
	defaultTimeout = 5
)

// execServer implements the GRPCExecServer interface
type execServer struct{}

func (s *execServer) ShowCmdTextOutput(a *pb.ShowCmdArgs, stream pb.GRPCExec_ShowCmdTextOutputServer) error {
	if a.GetCli() != defaultCmd {
		stream.Send(&pb.ShowCmdTextReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongCmdErr,
		})
		return errors.New(wrongCmdErr)
	}
	stream.Send(&pb.ShowCmdTextReply{
		ResReqId: a.GetReqId(),
		Output:   "show test output",
	})
	return nil
}

func (s *execServer) ShowCmdJSONOutput(a *pb.ShowCmdArgs, stream pb.GRPCExec_ShowCmdJSONOutputServer) error {
	if a.GetCli() != defaultCmd {
		stream.Send(&pb.ShowCmdJSONReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongCmdErr,
		})
		return errors.New(wrongCmdErr)
	}
	m := map[string]string{"result": "show test output"}
	j, err := json.Marshal(m)
	if err != nil {
		return errors.New("could not encode the test response")
	}
	stream.Send(&pb.ShowCmdJSONReply{
		ResReqId:   a.GetReqId(),
		Jsonoutput: string(j),
	})
	return nil
}

// operConfigServer implements the GRPCConfigOperServer interface
type operConfigServer struct{}

func (s *operConfigServer) GetConfig(a *pb.ConfigGetArgs, stream pb.GRPCConfigOper_GetConfigServer) error {
	if a.GetYangpathjson() != defaultYang {
		stream.Send(&pb.ConfigGetReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongYangErr,
		})
		return errors.New(wrongYangErr)
	}
	m := map[string]string{"result": "config"}
	j, err := json.Marshal(m)
	if err != nil {
		return errors.New("could not encode the test response")
	}
	stream.Send(&pb.ConfigGetReply{
		ResReqId: a.GetReqId(),
		Yangjson: string(j),
	})
	return nil
}

func (s *operConfigServer) MergeConfig(ctx context.Context, a *pb.ConfigArgs) (r *pb.ConfigReply, err error) {
	if a.GetYangjson() != defaultYang {
		err = errors.New(wrongYangErr)
		r = &pb.ConfigReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongYangErr,
		}
		return
	}
	r = &pb.ConfigReply{
		ResReqId: a.GetReqId(),
	}
	return
}

func (s *operConfigServer) DeleteConfig(ctx context.Context, a *pb.ConfigArgs) (r *pb.ConfigReply, err error) {
	if a.GetYangjson() != defaultYang {
		err = errors.New(wrongYangErr)
		r = &pb.ConfigReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongYangErr,
		}
		return
	}
	r = &pb.ConfigReply{
		ResReqId: a.GetReqId(),
	}
	return
}

func (s *operConfigServer) ReplaceConfig(ctx context.Context, a *pb.ConfigArgs) (r *pb.ConfigReply, err error) {
	if a.GetYangjson() != defaultYang {
		err = errors.New(wrongYangErr)
		r = &pb.ConfigReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongYangErr,
		}
		return
	}
	r = &pb.ConfigReply{
		ResReqId: a.GetReqId(),
	}
	return
}

func (s *operConfigServer) CliConfig(context.Context, *pb.CliConfigArgs) (*pb.CliConfigReply, error) {
	return nil, nil
}

func (s *operConfigServer) CommitReplace(context.Context, *pb.CommitReplaceArgs) (*pb.CommitReplaceReply, error) {
	return nil, nil
}

func (s *operConfigServer) CommitConfig(context.Context, *pb.CommitArgs) (*pb.CommitReply, error) {
	return nil, nil
}

func (s *operConfigServer) ConfigDiscardChanges(context.Context, *pb.DiscardChangesArgs) (*pb.DiscardChangesReply, error) {
	return nil, nil
}

func (s *operConfigServer) GetOper(*pb.GetOperArgs, pb.GRPCConfigOper_GetOperServer) error {
	return nil
}

func (s *operConfigServer) CreateSubs(*pb.CreateSubsArgs, pb.GRPCConfigOper_CreateSubsServer) error {
	return nil
}

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
		return errors.New(wrongCreds)
	}
	return errors.New("empty metadata")
}

func Server(t *testing.T, svc string) *grpc.Server {
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
	switch svc {
	case "exec":
		pb.RegisterGRPCExecServer(s, &execServer{})
	case "opercon":
		pb.RegisterGRPCConfigOperServer(s, &operConfigServer{})
	default:
	}

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

func TestConnect(t *testing.T) {
	x := xr.CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing.
		// Cert and Key for localhost are provided in the test folder
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}
	tt := []struct {
		name   string
		target string
		certf  string
		err    string
	}{
		{name: "local connection"},
		{name: "wrong target", target: "192.168.0.1:57344", err: "TBD"},
		{name: "wrong certificate", certf: "example/input/ems5502-1.pem", err: "TBD"},
		{name: "inexistent certificate", certf: "dummy", err: "TBD"},
	}
	s := Server(t, "none")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Get a copy of 'x' and change parameters of test case so requires
			xc := x
			if tc.certf != "" {
				xc.Cert = tc.certf
			}
			if tc.target != "" {
				xc.Host = tc.target
			}
			// Won't return Error if it keeps re-trying.
			conn, ctx, err := xr.Connect(xc)
			if err != nil {
				if tc.err != "" {
					return
				}
				t.Fatalf("could not setup a client connection to %v", xc.Host)
			}
			select {
			case <-ctx.Done():
				t.Fatalf("could not setup a client connection to %v in under 1.5 seconds", xc.Host)
			// Just wait for 2.5 seconds for this local connection to be setup.
			case <-time.After(2500 * time.Millisecond):
				break
			}
			// Connection won't fail until it timeouts. It re-attempt to connect until this happens.
			// Can initially timeout because of the WithTimeout option hard-coded to two seconds
			// or after an overal timeout of 'x.Timeout'
			conn.Close()
		})
	}
	s.Stop()
	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
	// reports 'bind: address already in use' when trying to run the next function test
	time.Sleep(200 * time.Millisecond)
}

func TestShowCmdTextOutput(t *testing.T) {
	x := xr.CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing.
		// Cert and Key for localhost are provided in the test folder
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}
	tt := []struct {
		name string
		cmd  string
		user string
		pass string
		err  string
	}{
		{name: "local connection", cmd: defaultCmd},
		{name: "wrong command", cmd: wrongCmd, err: wrongCmdErr},
		// TODO Fix the StreamInterceptor to hadle wrong authentication.
		// {name: "wrong user", cmd: defaultCmd, user: "bob", err: wrongCreds},
	}
	s := Server(t, "exec")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		if err != nil {
			t.Fatalf("could not setup a client connection to %v", x.Host)
		}

	}
	var id int64 = 1
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.user != "" {
				conn.Close()
				xc := x
				xc.User = tc.user
				conn, ctx, err = xr.Connect(xc)
				if err != nil {
					t.Fatalf("could not setup a client connection to %v", x.Host)
				}
			}
			_, err := xr.ShowCmdTextOutput(ctx, conn, tc.cmd, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongCmdErr) && tc.err == wrongCmdErr {
					return
				}
				if strings.Contains(err.Error(), wrongCreds) && tc.err == wrongCreds {
					return
				}
				t.Fatalf("failed to get show command text output from %v", x.Host)
			}
		})
		id++
	}
	conn.Close()
	s.Stop()
	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
	// reports 'bind: address already in use' when trying to run the next function test
	time.Sleep(200 * time.Millisecond)
}

func TestShowCmdJSONOutput(t *testing.T) {
	x := xr.CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing.
		// Cert and Key for localhost are provided in the test folder
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}
	tt := []struct {
		name string
		cmd  string
		err  string
	}{
		{name: "local connection", cmd: defaultCmd},
		{name: "wrong command", cmd: wrongCmd, err: wrongCmdErr},
	}
	s := Server(t, "exec")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		if err != nil {
			t.Fatalf("could not setup a client connection to %v", x.Host)
		}

	}
	var id int64 = 1
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := xr.ShowCmdJSONOutput(ctx, conn, tc.cmd, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongCmdErr) && tc.err == wrongCmdErr {
					return
				}
				t.Fatalf("failed to get show command json output from %v", x.Host)
			}
		})
		id++
	}
	conn.Close()
	s.Stop()
	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
	// reports 'bind: address already in use' when trying to run the next function test
	time.Sleep(200 * time.Millisecond)
}

func TestGetConfig(t *testing.T) {

}

func TestMergeConfig(t *testing.T) {
	x := xr.CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing.
		// Cert and Key for localhost are provided in the test folder
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}
	tt := []struct {
		name string
		conf string
		user string
		pass string
		err  string
	}{
		// The order of these test do matter, we change credentials
		// on the last ones.
		{name: "local connection", conf: defaultYang},
		{name: "wrong config", conf: "confreg 0x00", err: wrongYangErr},
		{name: "wrong user", conf: defaultYang, user: "bob", err: wrongCreds},
		{name: "wrong password", conf: defaultYang, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		if err != nil {
			t.Fatalf("could not setup a client connection to %v", x.Host)
		}
	}
	var id int64 = 1
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.user != "" {
				conn.Close()
				xc := x
				xc.User = tc.user
				conn, ctx, err = xr.Connect(xc)
				if err != nil {
					t.Fatalf("could not setup a client connection to %v", x.Host)
				}
			}
			if tc.pass != "" {
				conn.Close()
				xc := x
				xc.Password = tc.pass
				conn, ctx, err = xr.Connect(xc)
				if err != nil {
					t.Fatalf("could not setup a client connection to %v", x.Host)
				}
			}
			_, err := xr.MergeConfig(ctx, conn, tc.conf, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongCreds) && tc.err == wrongCreds {
					return
				}
				if strings.Contains(err.Error(), wrongYangErr) && tc.err == wrongYangErr {
					return
				}
				t.Fatalf("Incorrect response from %v, %v", x.Host, err)
			}
		})
		id++
	}
	conn.Close()
	s.Stop()
	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
	// reports 'bind: address already in use' when trying to run the next function test
	time.Sleep(200 * time.Millisecond)
}

func TestDeleteConfig(t *testing.T) {
	x := xr.CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing.
		// Cert and Key for localhost are provided in the test folder
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}
	tt := []struct {
		name string
		conf string
		user string
		pass string
		err  string
	}{
		// The order of these test do matter, we change credentials
		// on the last ones.
		{name: "local connection", conf: defaultYang},
		{name: "wrong config", conf: "confreg 0x00", err: wrongYangErr},
		{name: "wrong user", conf: defaultYang, user: "bob", err: wrongCreds},
		{name: "wrong password", conf: defaultYang, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		if err != nil {
			t.Fatalf("could not setup a client connection to %v", x.Host)
		}
	}
	var id int64 = 1
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.user != "" {
				conn.Close()
				xc := x
				xc.User = tc.user
				conn, ctx, err = xr.Connect(xc)
				if err != nil {
					t.Fatalf("could not setup a client connection to %v", x.Host)
				}
			}
			if tc.pass != "" {
				conn.Close()
				xc := x
				xc.Password = tc.pass
				conn, ctx, err = xr.Connect(xc)
				if err != nil {
					t.Fatalf("could not setup a client connection to %v", x.Host)
				}
			}
			_, err := xr.DeleteConfig(ctx, conn, tc.conf, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongCreds) && tc.err == wrongCreds {
					return
				}
				if strings.Contains(err.Error(), wrongYangErr) && tc.err == wrongYangErr {
					return
				}
				t.Fatalf("Incorrect response from %v, %v", x.Host, err)
			}
		})
		id++
	}
	conn.Close()
	s.Stop()
	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
	// reports 'bind: address already in use' when trying to run the next function test
	time.Sleep(200 * time.Millisecond)
}

func TestReplaceConfig(t *testing.T) {
	x := xr.CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing.
		// Cert and Key for localhost are provided in the test folder
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}
	tt := []struct {
		name string
		conf string
		user string
		pass string
		err  string
	}{
		// The order of these test do matter, we change credentials
		// on the last ones.
		{name: "local connection", conf: defaultYang},
		{name: "wrong config", conf: "confreg 0x00", err: wrongYangErr},
		{name: "wrong user", conf: defaultYang, user: "bob", err: wrongCreds},
		{name: "wrong password", conf: defaultYang, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		if err != nil {
			t.Fatalf("could not setup a client connection to %v", x.Host)
		}
	}
	var id int64 = 1
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.user != "" {
				conn.Close()
				xc := x
				xc.User = tc.user
				conn, ctx, err = xr.Connect(xc)
				if err != nil {
					t.Fatalf("could not setup a client connection to %v", x.Host)
				}
			}
			if tc.pass != "" {
				conn.Close()
				xc := x
				xc.Password = tc.pass
				conn, ctx, err = xr.Connect(xc)
				if err != nil {
					t.Fatalf("could not setup a client connection to %v", x.Host)
				}
			}
			_, err := xr.ReplaceConfig(ctx, conn, tc.conf, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongCreds) && tc.err == wrongCreds {
					return
				}
				if strings.Contains(err.Error(), wrongYangErr) && tc.err == wrongYangErr {
					return
				}
				t.Fatalf("Incorrect response from %v, %v", x.Host, err)
			}
		})
		id++
	}
	conn.Close()
	s.Stop()
	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
	// reports 'bind: address already in use' when trying to run the next function test
	time.Sleep(200 * time.Millisecond)
}
