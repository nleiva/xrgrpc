package xrgrpc

import (
	"fmt"
	"io"
	"net"
	"context"
	"errors"

	pb "github.com/nleiva/xrgrpc/proto/sla"
	"google.golang.org/grpc"
)

// VRFOperation handles VRF registration operations
// 	SL_REGOP_REGISTER    = 0x1 // VRF registration
// 	SL_REGOP_UNREGISTER  = 0x2 // VRF Un-registeration
// 	SL_REGOP_EOF         = 0x3 // After Registration, the client is expected to send an EOF
func VRFOperation(conn *grpc.ClientConn, o int, d uint32) error {
	// These are two hard-coded variables. TODO; take them as input
	v := "default"
	var p uint32 = 500

	op := new(pb.SLRegOp)
	switch o {
	case 1:
		*op = pb.SLRegOp_SL_REGOP_REGISTER
	case 2:
		*op = pb.SLRegOp_SL_REGOP_UNREGISTER
	case 3:
		*op = pb.SLRegOp_SL_REGOP_EOF
	default:
		return errors.New("Unidentified VRF Operation")
	}

	// 'c' is the gRPC stub.
	c := pb.NewSLRoutev6OperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.SLVrfRegMsg{
		Oper: *op,
		VrfRegMsgs: []*pb.SLVrfReg{
			&pb.SLVrfReg{
				VrfName:                 v,
				AdminDistance:           d,
				VrfPurgeIntervalSeconds: p,
			},
		},
	}

	// 'r' is the result that comes back from the target.
	r, err := c.SLRoutev6VrfRegOp(context.Background(), &a)
	if err != nil {
		return fmt.Errorf("gRPC SLRoutev6VrfRegOp failed: %w", err)
	}

	//	SL_SUCCESS: Entire bulk operation was successful.
	//	SL_SOME_ERR: Operation failed for one or more entries.
	//	SL_RPC_XXX: Eentire bulk operation failed.
	if r.GetStatusSummary().GetStatus() != pb.SLErrorStatus_SL_SUCCESS {
		// TODO: Add cases for partial errors
		return fmt.Errorf("Error triggered by remote host: %s", r.GetStatusSummary().GetStatus().String())
	}
	return nil
}

// SetRoute ...
// 	SL_OBJOP_ADD = 1 // Route add. Fails if the route already exists.
// 	SL_OBJOP_UPDATE = 2 // Route update. Creates or updates the route.
// 	SL_OBJOP_DELETE = 3 // Route delete. The route path is not necessary to delete the route.
//	SLRoutev6Op(ctx context.Context, in *SLRoutev6Msg, opts ...grpc.CallOption) (*SLRoutev6MsgRsp, error)
func SetRoute(conn *grpc.ClientConn, o int, ad string, d uint32, nh string) error {
	// These are two hard-coded variables. TODO; take them as input
	v := "default"
	//intf := "HundredGigE0/0/0/1"
	_, nw, err := net.ParseCIDR(ad)
	if err != nil {
		return fmt.Errorf("Could not parse address: %w", err)
	}
	mk, _ := nw.Mask.Size()

	op := new(pb.SLObjectOp)
	switch o {
	case 1:
		*op = pb.SLObjectOp_SL_OBJOP_ADD
	case 2:
		*op = pb.SLObjectOp_SL_OBJOP_UPDATE
	case 3:
		*op = pb.SLObjectOp_SL_OBJOP_DELETE
	default:
		return errors.New("Unidentified Object Operation")
	}

	// 'c' is the gRPC stub.
	c := pb.NewSLRoutev6OperClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.SLRoutev6Msg{
		Oper: *op,
		// Correlator: x,
		VrfName: v,
		Routes: []*pb.SLRoutev6{
			&pb.SLRoutev6{
				Prefix:    nw.IP,
				PrefixLen: uint32(mk),
				RouteCommon: &pb.SLRouteCommon{
					AdminDistance: d,
				},
				PathList: []*pb.SLRoutePath{
					&pb.SLRoutePath{
						// NexthopInterface: x,
						NexthopAddress: &pb.SLIpAddress{
							Address: &pb.SLIpAddress_V6Address{
								V6Address: net.ParseIP(nh),
							},
						},
					},
				},
			},
		},
	}

	// 'r' is the result that comes back from the target.
	r, err := c.SLRoutev6Op(context.Background(), &a)
	if err != nil {
		return fmt.Errorf("gRPC SLRoutev6Op failed: %w", err)
	}

	//	SL_SUCCESS: Entire bulk operation was successful.
	//	SL_SOME_ERR: Operation failed for one or more entries.
	//	SL_RPC_XXX: Eentire bulk operation failed.
	if r.GetStatusSummary().GetStatus() != pb.SLErrorStatus_SL_SUCCESS {
		// TODO: Add cases for partial errors
		return fmt.Errorf("Error triggered by remote host: %s, AND %v",
			r.GetStatusSummary().GetStatus().String(), r.GetResults()[0].GetErrStatus().GetStatus())
	}
	return nil
}

// ClientInit takes care of global initialization and setup
// a notification channel
func ClientInit(conn *grpc.ClientConn) error {
	ch := make(chan int)
	/* Setup the notification channel */
	go setupNotifChannel(conn, ch)

	/* Wait for response 0: error. 1: all ok*/
	if w := <-ch; w == 0 {
		return errors.New("Error triggered by remote host")
	}

	// 'c' is the gRPC stub.
	c := pb.NewSLGlobalClient(conn)

	// 'a' is the object we send to the router via the stub.
	a := pb.SLGlobalsGetMsg{}

	// 'r' is the result that comes back from the target.
	_, err := c.SLGlobalsGet(context.Background(), &a)
	if err != nil {
		return fmt.Errorf("gRPC SLGlobalsGet failed: %w", err)
	}
	return nil
}

func setupNotifChannel(conn *grpc.ClientConn, ch chan int) {
	// 'c' is the gRPC stub.
	c := pb.NewSLGlobalClient(conn)

	// 'a' is the object we send to the router via the stub.
	// Version Major.Minor.Subversion define the API's version.
	// Hardcoded to 0.0.1 for now
	a := pb.SLInitMsg{
		MajorVer: uint32(0),
		MinorVer: uint32(0),
		SubVer:   uint32(1),
	}

	// 'st' is the streamed result that comes back from the target.
	st, err := c.SLGlobalInitNotif(context.Background(), &a)
	if err != nil {
		ch <- 0
		return
		// return errors.Wrap(err, "gRPC SLGlobalInitNotif failed")
	}

	for {
		// Loop through the responses in the stream until there is nothing left.
		r, err := st.Recv()
		if err == io.EOF {
			return
			//return nil
		}
		if err != nil {
			return
			// return errors.Wrap(err, "Recive Stream failed")
		}

		// Status code, interpreted based on the Event Type.
		//
		//   case EventType == SL_GLOBAL_EVENT_TYPE_ERROR:
		//       case ErrStatus == SL_NOTIF_TERM:
		//          => Another client is attempting to take over the session.
		//             This session will be closed.
		//       case ErrStatus == (some error from SLErrorStatus)
		//          => Client must look into the specific error message returned.
		//
		//   case EventType == SL_GLOBAL_EVENT_TYPE_HEARTBEAT:
		//       case ErrStatus == SL_SUCCESS
		//          => Client can safely ignore this heartbeat message.
		//
		//   case EventType == SL_GLOBAL_EVENT_TYPE_VERSION:
		//       case ErrStatus == SL_SUCCESS
		//          => Client version accepted.
		//       case ErrStatus == SL_INIT_STATE_READY
		//          => Client version accepted.
		//             Any previous state was sucessfully recovered.
		//       case ErrStatus == SL_INIT_STATE_CLEAR
		//          => Client version accepted. Any previous state was lost.
		//             Client must replay all previous objects to server.
		//       case ErrStatus == SL_UNSUPPORTED_VER
		//          => Client and Server version mismatch. The client is not
		//             allowed to proceed, and the channel will be closed.
		//       case ErrStatus == (some error from SLErrorStatus)
		//          => Client must either try again, or look into the specific
		//             error message returned.

		// Need to Fix this!!
		switch r.EventType {
		case pb.SLGlobalNotifType_SL_GLOBAL_EVENT_TYPE_VERSION:
			ch <- 1
			continue
		case pb.SLGlobalNotifType_SL_GLOBAL_EVENT_TYPE_ERROR:
			_ = fmt.Errorf("%s", r.ErrStatus.String())
			break
		case pb.SLGlobalNotifType_SL_GLOBAL_EVENT_TYPE_HEARTBEAT:
			continue
		default:
			_ = fmt.Errorf("%s", r.ErrStatus.String())
			return
		}
	}
}
