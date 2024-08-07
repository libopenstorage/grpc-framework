// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: apis/example/apiv1/example.proto

/*
Package example is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package example

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_ExampleGreeter_SayExample_0(ctx context.Context, marshaler runtime.Marshaler, client ExampleGreeterClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ExampleGreeterSayExampleRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.SayExample(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ExampleGreeter_SayExample_0(ctx context.Context, marshaler runtime.Marshaler, server ExampleGreeterServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ExampleGreeterSayExampleRequest
	var metadata runtime.ServerMetadata

	if err := marshaler.NewDecoder(req.Body).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.SayExample(ctx, &protoReq)
	return msg, metadata, err

}

func request_ExampleIdentity_ServerVersion_0(ctx context.Context, marshaler runtime.Marshaler, client ExampleIdentityClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ExampleIdentityVersionRequest
	var metadata runtime.ServerMetadata

	msg, err := client.ServerVersion(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ExampleIdentity_ServerVersion_0(ctx context.Context, marshaler runtime.Marshaler, server ExampleIdentityServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ExampleIdentityVersionRequest
	var metadata runtime.ServerMetadata

	msg, err := server.ServerVersion(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterExampleGreeterHandlerServer registers the http handlers for service ExampleGreeter to "mux".
// UnaryRPC     :call ExampleGreeterServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterExampleGreeterHandlerFromEndpoint instead.
func RegisterExampleGreeterHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ExampleGreeterServer) error {

	mux.Handle("POST", pattern_ExampleGreeter_SayExample_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/hello.example.v1.ExampleGreeter/SayExample", runtime.WithHTTPPathPattern("/v1/greeter:sayExample"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ExampleGreeter_SayExample_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ExampleGreeter_SayExample_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterExampleIdentityHandlerServer registers the http handlers for service ExampleIdentity to "mux".
// UnaryRPC     :call ExampleIdentityServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterExampleIdentityHandlerFromEndpoint instead.
func RegisterExampleIdentityHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ExampleIdentityServer) error {

	mux.Handle("GET", pattern_ExampleIdentity_ServerVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/hello.example.v1.ExampleIdentity/ServerVersion", runtime.WithHTTPPathPattern("/v1/identity:serverVersion"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ExampleIdentity_ServerVersion_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ExampleIdentity_ServerVersion_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterExampleGreeterHandlerFromEndpoint is same as RegisterExampleGreeterHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterExampleGreeterHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterExampleGreeterHandler(ctx, mux, conn)
}

// RegisterExampleGreeterHandler registers the http handlers for service ExampleGreeter to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterExampleGreeterHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterExampleGreeterHandlerClient(ctx, mux, NewExampleGreeterClient(conn))
}

// RegisterExampleGreeterHandlerClient registers the http handlers for service ExampleGreeter
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ExampleGreeterClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ExampleGreeterClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ExampleGreeterClient" to call the correct interceptors.
func RegisterExampleGreeterHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ExampleGreeterClient) error {

	mux.Handle("POST", pattern_ExampleGreeter_SayExample_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/hello.example.v1.ExampleGreeter/SayExample", runtime.WithHTTPPathPattern("/v1/greeter:sayExample"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ExampleGreeter_SayExample_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ExampleGreeter_SayExample_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ExampleGreeter_SayExample_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "greeter"}, "sayExample"))
)

var (
	forward_ExampleGreeter_SayExample_0 = runtime.ForwardResponseMessage
)

// RegisterExampleIdentityHandlerFromEndpoint is same as RegisterExampleIdentityHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterExampleIdentityHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterExampleIdentityHandler(ctx, mux, conn)
}

// RegisterExampleIdentityHandler registers the http handlers for service ExampleIdentity to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterExampleIdentityHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterExampleIdentityHandlerClient(ctx, mux, NewExampleIdentityClient(conn))
}

// RegisterExampleIdentityHandlerClient registers the http handlers for service ExampleIdentity
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ExampleIdentityClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ExampleIdentityClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ExampleIdentityClient" to call the correct interceptors.
func RegisterExampleIdentityHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ExampleIdentityClient) error {

	mux.Handle("GET", pattern_ExampleIdentity_ServerVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/hello.example.v1.ExampleIdentity/ServerVersion", runtime.WithHTTPPathPattern("/v1/identity:serverVersion"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ExampleIdentity_ServerVersion_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ExampleIdentity_ServerVersion_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ExampleIdentity_ServerVersion_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "identity"}, "serverVersion"))
)

var (
	forward_ExampleIdentity_ServerVersion_0 = runtime.ForwardResponseMessage
)
