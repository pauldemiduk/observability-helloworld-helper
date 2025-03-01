// load_context.go - Optimized structured logging for the Hello World microservice
package models

// LoadLogContext -
// Initializes and returns a LogContext struct.
// It holds key metadata about the running service instance which is captured/written in log attributes
func LoadLogContext() LogContext {
	return LogContext{
		ServiceName:      HelloInstance.Context.ServiceName,
		ServiceVersion:   HelloInstance.Context.ServiceVersion,
		ServiceHostname:  HelloInstance.Context.ServiceHostname,
		ServiceIP:        HelloInstance.Context.ServiceIP,
		ServicePort:      HelloInstance.Context.ServicePort,
		ProcessID:        HelloInstance.Context.ProcessID,
		ThreadID:         HelloInstance.Context.ThreadID,
		NodeIP:           HelloInstance.Context.NodeIP,
		NodePort:         HelloInstance.Context.NodePort,
		CloudProvider:    HelloInstance.Context.CloudProvider,
		CloudRegion:      HelloInstance.Context.CloudRegion,
		AvailabilityZone: HelloInstance.Context.AvailabilityZone,
		K8sNamespace:     HelloInstance.Context.K8sNamespace,
		K8sPodName:       HelloInstance.Context.K8sPodName,
		K8sContainerID:   HelloInstance.Context.K8sContainerID,
		LoadBalancerIP:   HelloInstance.Context.LoadBalancerIP,
		InstanceID:       HelloInstance.Context.InstanceID,
	}
}
