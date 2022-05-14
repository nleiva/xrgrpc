package xrgrpc_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	xr "github.com/nleiva/xrgrpc"
	pb "github.com/nleiva/xrgrpc/proto/ems"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	defaultAddr            = "localhost"
	defaultPort            = ":57344"
	defaultUser            = "test"
	defaultPass            = "test"
	defaultCert            = "test/cert.pem"
	defaultKey             = "test/key.pem"
	defaultCmd             = "show test"
	defaultYang            = "{\"Cisco-IOS-XR-test:tree\": [null]}"
	defaultSubsID          = "TEST"
	defaultCommitID uint32 = 100000002
	wrongCmd               = "show me the money"
	wrongConf              = "confreg 0x00"
	wrongYang              = "{\"Cisco-IOS-XR-fake:tree\": [null]}"
	wrongCreds             = "incorrect username/password"
	wrongSubsID            = "wrong Subscription ID"
	wrongEncode            = "wrong encoding"
	wrongCmdErr            = "wrong command"
	wrongYangErr           = "wrong YANG path"
	wrongCommitID          = "wrong Commit ID"
	defaultTimeout         = 5
)

// execServer implements the GRPCExecServer interface
type execServer struct{ pb.UnimplementedGRPCExecServer }

func (s *execServer) ShowCmdTextOutput(a *pb.ShowCmdArgs, stream pb.GRPCExec_ShowCmdTextOutputServer) error {
	if a.GetCli() != defaultCmd {
		err := stream.Send(&pb.ShowCmdTextReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongCmdErr,
		})
		if err != nil {
			return err
		}
		return errors.New(wrongCmdErr)
	}
	err := stream.Send(&pb.ShowCmdTextReply{
		ResReqId: a.GetReqId(),
		Output:   "show test output",
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *execServer) ShowCmdJSONOutput(a *pb.ShowCmdArgs, stream pb.GRPCExec_ShowCmdJSONOutputServer) error {
	if a.GetCli() != defaultCmd {
		err := stream.Send(&pb.ShowCmdJSONReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongCmdErr,
		})
		if err != nil {
			return err
		}
		return errors.New(wrongCmdErr)
	}
	m := map[string]string{"result": "show test output"}
	j, err := json.Marshal(m)
	if err != nil {
		return errors.New("could not encode the test response")
	}
	err = stream.Send(&pb.ShowCmdJSONReply{
		ResReqId:   a.GetReqId(),
		Jsonoutput: string(j),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *execServer) ActionJSON(a *pb.ActionJSONArgs, stream pb.GRPCExec_ActionJSONServer) error {
	if a.GetYangpathjson() != defaultYang {
		err := stream.Send(&pb.ActionJSONReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongCmdErr,
		})
		if err != nil {
			return err
		}
		return errors.New(wrongCmdErr)
	}
	m := map[string]string{"result": "action test output"}
	j, err := json.Marshal(m)
	if err != nil {
		return errors.New("could not encode the test response")
	}
	err = stream.Send(&pb.ActionJSONReply{
		ResReqId: a.GetReqId(),
		Yangjson: string(j),
	})
	if err != nil {
		return err
	}
	return nil
}

// operConfigServer implements the GRPCConfigOperServer interface
type operConfigServer struct {
	pb.UnimplementedGRPCConfigOperServer
}

func (s *operConfigServer) GetConfig(a *pb.ConfigGetArgs, stream pb.GRPCConfigOper_GetConfigServer) error {
	if a.GetYangpathjson() != defaultYang {
		err := stream.Send(&pb.ConfigGetReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongYangErr,
		})
		if err != nil {
			return err
		}
		return errors.New(wrongYangErr)
	}
	m := map[string]string{"result": "config"}
	j, err := json.Marshal(m)
	if err != nil {
		return errors.New("could not encode the test response")
	}
	err = stream.Send(&pb.ConfigGetReply{
		ResReqId: a.GetReqId(),
		Yangjson: string(j),
	})
	if err != nil {
		return err
	}
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

func (s *operConfigServer) CliConfig(ctx context.Context, a *pb.CliConfigArgs) (r *pb.CliConfigReply, err error) {
	if a.GetCli() != defaultCmd {
		err = errors.New(wrongCmdErr)
		r = &pb.CliConfigReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongCmdErr,
		}
		return
	}
	r = &pb.CliConfigReply{
		ResReqId: a.GetReqId(),
	}
	return
}

// TODO: Add test case for this!
func (s *operConfigServer) CommitReplace(ctx context.Context, a *pb.CommitReplaceArgs) (r *pb.CommitReplaceReply, err error) {
	return
}

// CommitConfig commits a config. Need to clarify its use-case.
func (s *operConfigServer) CommitConfig(ctx context.Context, a *pb.CommitArgs) (r *pb.CommitReply, err error) {
	if a.GetCommitID() != defaultCommitID {
		err = errors.New(wrongCommitID)
		r = &pb.CommitReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongCommitID,
		}
		return
	}
	r = &pb.CommitReply{
		ResReqId: a.GetReqId(),
	}
	return
}

// CommitConfig commits a config. Need to clarify its use-case.
func (s *operConfigServer) ConfigDiscardChanges(context.Context, *pb.DiscardChangesArgs) (*pb.DiscardChangesReply, error) {
	return nil, nil
}

func (s *operConfigServer) GetOper(a *pb.GetOperArgs, stream pb.GRPCConfigOper_GetOperServer) error {
	if a.GetYangpathjson() != defaultYang {
		err := stream.Send(&pb.GetOperReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongYangErr,
		})
		if err != nil {
			return err
		}
		return errors.New(wrongYangErr)
	}
	m := map[string]string{"result": "oper"}
	j, err := json.Marshal(m)
	if err != nil {
		return errors.New("could not encode the test response")
	}
	err = stream.Send(&pb.GetOperReply{
		ResReqId: a.GetReqId(),
		Yangjson: string(j),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *operConfigServer) CreateSubs(a *pb.CreateSubsArgs, stream pb.GRPCConfigOper_CreateSubsServer) error {
	mape := map[int64]string{
		2: "gpb",
		3: "gpbkv",
		4: "json",
	}
	_, ok := mape[a.GetEncode()]
	if !ok {
		return fmt.Errorf("%s, '%v' not supported", wrongEncode, a.GetEncode())
	}
	if a.GetSubidstr() != defaultSubsID {
		err := stream.Send(&pb.CreateSubsReply{
			ResReqId: a.GetReqId(),
			Errors:   wrongSubsID,
		})
		if err != nil {
			return err
		}
		return errors.New(wrongSubsID)
	}
	m := map[string]string{"result": "oper"}
	j, err := json.Marshal(m)
	if err != nil {
		return errors.New("could not encode the test response")
	}
	// Telemetry fixed at 0.45 second interval for testing
	ticker := time.NewTicker(450 * time.Millisecond)

	// With this ('n') we can simulate server and client connection cancellation:
	// 	n < ctx.Timeout -> Server cancels
	// 	n > ctx.Timeout -> Client cancels
	// Considering the only numerical inputs we have are the ID and Encoding, we will
	// re-use the latter to timeout the stream.
	n := time.Duration(a.GetEncode()) * time.Second
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(n)
		timeout <- true
	}()
	for {
		select {
		case <-ticker.C:
			err := stream.Send(&pb.CreateSubsReply{
				ResReqId: a.GetReqId(),
				Data:     j,
			})
			if err != nil {
				return err
			}
		case <-timeout:
			ticker.Stop()
			return nil

		}
	}
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
	if md, ok := metadata.FromIncomingContext(ctx); ok {
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
		t.Fatalf("failed to construct TLS credentials: %v", err)
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
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			t.Fatalf("failed to serve: %v", err)
		}
	}()
	return s
}

func ServerInsecure(t *testing.T, svc string) *grpc.Server {
	lis, err := net.Listen("tcp", defaultPort)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	// var opts []grpc.ServerOption
	s := grpc.NewServer(
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
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
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
		{name: "wrong certificate", certf: "example/input/certificate/ems5502-1.pem", err: "TBD"},
		{name: "inexistent certificate", certf: "dummy", err: "TBD"},
		{name: "No certificate", certf: "empty"},
	}
	s := Server(t, "none")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Get a copy of 'x' and change parameters of test case so requires
			xc := x
			if tc.certf != "" {
				xc.Cert = tc.certf
			}
			if tc.certf == "empty" {
				xc.Cert = ""
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

func TestConnectInsecure(t *testing.T) {
	x := xr.CiscoGrpcClient{
		// User/Password are per RPC based, won't be checked when dialing.
		// Cert and Key for localhost are provided in the test folder
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
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
	}
	s := ServerInsecure(t, "none")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Get a copy of 'x' and change parameters of test case so requires
			xc := x
			// Won't return Error if it keeps re-trying.
			conn, ctx, err := xr.ConnectInsecure(xc)
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
		t.Fatalf("could not setup a client connection to %v", x.Host)
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
		t.Fatalf("could not setup a client connection to %v", x.Host)
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

// func TestActionJSONOutput(t *testing.T) {
// 	x := xr.CiscoGrpcClient{
// 		User:     defaultUser,
// 		Password: defaultPass,
// 		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
// 		Cert:     defaultCert,
// 		Domain:   "localhost",
// 		Timeout:  defaultTimeout,
// 	}
// 	tt := []struct {
// 		name string
// 		act  string
// 		err  string
// 	}{
// 		{name: "local connection", act: defaultYang},
// 		{name: "wrong command", act: wrongCmd, err: wrongCmdErr},
// 	}
// 	s := Server(t, "exec")
// 	conn, ctx, err := xr.Connect(x)
// 	if err != nil {
// 		t.Fatalf("could not setup a client connection to %v", x.Host)
// 	}
// 	var id int64 = 1
// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			_, err := xr.ActionJSON(ctx, conn, tc.act, id)
// 			if err != nil {
// 				if strings.Contains(err.Error(), wrongCmdErr) && tc.err == wrongCmdErr {
// 					return
// 				}
// 				t.Fatalf("failed to get action json output from %v", x.Host)
// 			}
// 		})
// 		id++
// 	}
// 	conn.Close()
// 	s.Stop()
// 	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
// 	// reports 'bind: address already in use' when trying to run the next function test
// 	time.Sleep(200 * time.Millisecond)
// }

func TestGetConfig(t *testing.T) {
	x := xr.CiscoGrpcClient{
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}

	tt := []struct {
		name  string
		paths string
		enc   int64
		err   string
	}{
		{name: "local connection", paths: defaultYang},
		{name: "wrong paths", paths: wrongYang, err: wrongYangErr},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		t.Fatalf("could not setup a client connection to %v", x.Host)
	}
	var id int64 = 1
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := xr.GetConfig(ctx, conn, tc.paths, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongYangErr) && tc.err == wrongYangErr {
					return
				}
				t.Fatalf("failed to get config from %v", x.Host)
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

func TestMergeConfig(t *testing.T) {
	x := xr.CiscoGrpcClient{
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
		{name: "wrong config", conf: wrongYang, err: wrongYangErr},
		{name: "wrong user", conf: defaultYang, user: "bob", err: wrongCreds},
		{name: "wrong password", conf: defaultYang, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		t.Fatalf("could not setup a client connection to %v", x.Host)
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
				t.Fatalf("incorrect response from %v, %v", x.Host, err)
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
		{name: "wrong config", conf: wrongYang, err: wrongYangErr},
		{name: "wrong user", conf: defaultYang, user: "bob", err: wrongCreds},
		{name: "wrong password", conf: defaultYang, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		t.Fatalf("could not setup a client connection to %v", x.Host)
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
				t.Fatalf("incorrect response from %v, %v", x.Host, err)
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
		{name: "wrong config", conf: wrongYang, err: wrongYangErr},
		{name: "wrong user", conf: defaultYang, user: "bob", err: wrongCreds},
		{name: "wrong password", conf: defaultYang, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		t.Fatalf("could not setup a client connection to %v", x.Host)
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
				t.Fatalf("incorrect response from %v, %v", x.Host, err)
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

func TestCLIConfig(t *testing.T) {
	x := xr.CiscoGrpcClient{
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
		{name: "local connection", conf: defaultCmd},
		{name: "wrong config", conf: wrongConf, err: wrongCmdErr},
		{name: "wrong user", conf: defaultCmd, user: "bob", err: wrongCreds},
		{name: "wrong password", conf: defaultCmd, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		t.Fatalf("could not setup a client connection to %v", x.Host)
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
			err := xr.CLIConfig(ctx, conn, tc.conf, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongCreds) && tc.err == wrongCreds {
					return
				}
				if strings.Contains(err.Error(), wrongCmdErr) && tc.err == wrongCmdErr {
					return
				}
				t.Fatalf("incorrect response from %v, %v", x.Host, err)
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

func TestCommitConfig(t *testing.T) {
	x := xr.CiscoGrpcClient{
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		Timeout:  defaultTimeout,
	}

	tt := []struct {
		name string
		cid  uint32
		user string
		pass string
		err  string
	}{
		// The order of these test do matter, we change credentials
		// on the last ones.
		{name: "local connection", cid: defaultCommitID},
		{name: "wrong config", cid: defaultCommitID, err: wrongCmdErr},
		{name: "wrong user", cid: defaultCommitID, user: "bob", err: wrongCreds},
		{name: "wrong password", cid: defaultCommitID, pass: "password", err: wrongCreds},
	}
	s := Server(t, "opercon")
	conn, ctx, err := xr.Connect(x)
	if err != nil {
		t.Fatalf("could not setup a client connection to %v", x.Host)
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
			_, err := xr.CommitConfig(ctx, conn, tc.cid, id)
			if err != nil {
				if strings.Contains(err.Error(), wrongCreds) && tc.err == wrongCreds {
					return
				}
				if strings.Contains(err.Error(), wrongCmdErr) && tc.err == wrongCmdErr {
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

func TestGetSubscription(t *testing.T) {
	x := xr.CiscoGrpcClient{
		User:     defaultUser,
		Password: defaultPass,
		Host:     strings.Join([]string{defaultAddr, defaultPort}, ""),
		Cert:     defaultCert,
		Domain:   "localhost",
		// We fixed Timeout to 3, in this case, in order to test different failure scenarios
		Timeout: 3,
	}

	tt := []struct {
		name string
		subs string
		enc  int64
		err  string
	}{
		{name: "server timeout", subs: defaultSubsID, enc: 2},
		{name: "client timeout", subs: defaultSubsID, enc: 4},
		{name: "wrong subscription", subs: "anything", enc: 3, err: wrongSubsID},
		{name: "wrong encoding", subs: defaultSubsID, enc: 5, err: wrongEncode},
	}
	s := Server(t, "opercon")
	var id int64 = 1

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// start := time.Now()
			conn, ctx, err := xr.Connect(x)
			defer conn.Close()
			if err != nil {
				t.Fatalf("could not setup a client connection to %v", x.Host)
			}
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			ch, ech, err := xr.GetSubscription(ctx, conn, tc.subs, id, tc.enc)
			if err != nil {
				if strings.Contains(err.Error(), wrongSubsID) && tc.err == wrongSubsID ||
					strings.Contains(err.Error(), wrongEncode) && tc.err == wrongEncode {
					return
				}
				t.Fatalf("could not setup Telemetry Subscription from %v: %v", x.Host, err)
			}
			// copy tc.err to avoid race condition
			go func(e string) {
				select {
				case <-ctx.Done():
					// Timeout: "context deadline exceeded"
					// err = ctx.Err()
					// fmt.Printf("\ngRPC session timed out after %v seconds: %v\n\n", time.Since(start), err.Error())
					return
				case err = <-ech:
					if err.Error() == "EOF" ||
						strings.Contains(err.Error(), wrongSubsID) && e == wrongSubsID ||
						strings.Contains(err.Error(), wrongEncode) && e == wrongEncode {
						return
					}
					// Session canceled: "context canceled"
					t.Fatalf("\ngRPC session to %v failed: %v\n\n", x.Host, err.Error())
				}
			}(tc.err)
			i := 1
			for tele := range ch {
				fmt.Printf("Telemetry Message %v-%v: %s\n", tc.enc, i, string(tele))
				i++
			}
			id++
		})

	}
	s.Stop()
	// To avoid tests failing in Travis CI, we sleep for 0.2 seconds, otherwise it
	// reports 'bind: address already in use' when trying to run the next function test
	time.Sleep(200 * time.Millisecond)
}

func TestBuildRouter(t *testing.T) {
	tt := []struct {
		name    string
		user    string
		pass    string
		host    string
		cert    string
		timeout int
		err     string
	}{
		{name: "default case", user: defaultUser, pass: defaultPass, host: defaultAddr + defaultPort, cert: defaultCert, timeout: defaultTimeout},
		{name: "wrong username", pass: defaultPass, host: defaultAddr + defaultPort, cert: defaultCert,
			timeout: defaultTimeout, err: "invalid username"},
		{name: "wrong password", user: defaultUser, host: defaultAddr + defaultPort, cert: defaultCert,
			timeout: defaultTimeout, err: "invalid password"},
		{name: "wrong host", user: defaultUser, pass: defaultPass, host: "300.1.1.1:57344", cert: defaultCert,
			timeout: defaultTimeout, err: "not a valid host address"},
		{name: "wrong cert file", user: defaultUser, pass: defaultPass, host: defaultAddr + defaultPort,
			timeout: defaultTimeout, err: "not a valid file location"},
		{name: "wrong timeout", user: defaultUser, pass: defaultPass, host: defaultAddr + defaultPort, cert: defaultCert,
			timeout: 0, err: "timeout must be greater than zero"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := xr.BuildRouter(
				xr.WithUsername(tc.user),
				xr.WithPassword(tc.pass),
				xr.WithHost(tc.host),
				xr.WithCert(tc.cert),
				xr.WithTimeout(tc.timeout),
			)
			if err != nil {
				if strings.Contains(err.Error(), tc.err) {
					return
				}
				t.Fatalf("Target parameters are incorrect: %s", err.Error())
			}

		})
	}

}
