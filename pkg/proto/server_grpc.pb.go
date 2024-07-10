// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.2
// source: pkg/proto/server.proto

package server_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Metric_UpdateMetrics_FullMethodName = "/server.Metric/UpdateMetrics"
)

// MetricClient is the client API for Metric service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricClient interface {
	UpdateMetrics(ctx context.Context, in *UpdateMetricsRequest, opts ...grpc.CallOption) (*UpdateMetricsResponse, error)
}

type metricClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricClient(cc grpc.ClientConnInterface) MetricClient {
	return &metricClient{cc}
}

func (c *metricClient) UpdateMetrics(ctx context.Context, in *UpdateMetricsRequest, opts ...grpc.CallOption) (*UpdateMetricsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetricsResponse)
	err := c.cc.Invoke(ctx, Metric_UpdateMetrics_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricServer is the server API for Metric service.
// All implementations must embed UnimplementedMetricServer
// for forward compatibility
type MetricServer interface {
	UpdateMetrics(context.Context, *UpdateMetricsRequest) (*UpdateMetricsResponse, error)
	mustEmbedUnimplementedMetricServer()
}

// UnimplementedMetricServer must be embedded to have forward compatible implementations.
type UnimplementedMetricServer struct {
}

func (UnimplementedMetricServer) UpdateMetrics(context.Context, *UpdateMetricsRequest) (*UpdateMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetrics not implemented")
}
func (UnimplementedMetricServer) mustEmbedUnimplementedMetricServer() {}

// UnsafeMetricServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricServer will
// result in compilation errors.
type UnsafeMetricServer interface {
	mustEmbedUnimplementedMetricServer()
}

func RegisterMetricServer(s grpc.ServiceRegistrar, srv MetricServer) {
	s.RegisterService(&Metric_ServiceDesc, srv)
}

func _Metric_UpdateMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServer).UpdateMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Metric_UpdateMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServer).UpdateMetrics(ctx, req.(*UpdateMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Metric_ServiceDesc is the grpc.ServiceDesc for Metric service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metric_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "server.Metric",
	HandlerType: (*MetricServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateMetrics",
			Handler:    _Metric_UpdateMetrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/proto/server.proto",
}
