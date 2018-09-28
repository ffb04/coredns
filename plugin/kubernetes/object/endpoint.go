package object

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Endpoints is a stripped down api.Endpoints with only the items we need for CoreDNS.
type Endpoints struct {
	Version   string
	Name      string
	Namespace string
	Index     string
	Subsets   []EndpointSubset

	*Empty
}

// EndpointSubset is a group of addresses with a common set of ports. The
// expanded set of endpoints is the Cartesian product of Addresses x Ports.
type EndpointSubset struct {
	Addresses []EndpointAddress
	Ports     []EndpointPort
}

// EndpointAddress is a tuple that describes single IP address.
type EndpointAddress struct {
	IP            string
	Hostname      string
	NodeName      string
	TargetRefName string
}

// EndpointPort is a tuple that describes a single port.
type EndpointPort struct {
	Port     int32
	Name     string
	Protocol string
}

// EndpointsKey return a string using for the index.
func EndpointsKey(name, namespace string) string { return name + "." + namespace }

// ToEndpoints converts an api.Service to a *Service.
func ToEndpoints(obj interface{}) interface{} {
	end, ok := obj.(*api.Endpoints)
	if !ok {
		return nil
	}

	e := &Endpoints{
		Version:   end.GetResourceVersion(),
		Name:      end.GetName(),
		Namespace: end.GetNamespace(),
		Index:     EndpointsKey(end.GetName(), end.GetNamespace()),
	}
	for _, eps := range end.Subsets {
		sub := EndpointSubset{}
		for _, a := range eps.Addresses {
			ea := EndpointAddress{IP: a.IP, Hostname: a.Hostname}
			if a.NodeName != nil {
				ea.NodeName = *a.NodeName
			}
			if a.TargetRef != nil {
				ea.TargetRefName = a.TargetRef.Name
			}
			sub.Addresses = append(sub.Addresses, ea)
		}
		for _, p := range eps.Ports {
			ep := EndpointPort{Port: p.Port, Name: p.Name, Protocol: string(p.Protocol)}
			sub.Ports = append(sub.Ports, ep)
		}
		// Add sentinal is there are no ports.
		if len(eps.Ports) == 0 {
			sub.Ports = []EndpointPort{{Port: -1}}
		}
		e.Subsets = append(e.Subsets, sub)
	}

	return e
}

var _ runtime.Object = &Endpoints{}

// GetNamespace implements the metav1.Object interface.
func (e *Endpoints) GetNamespace() string { return e.Namespace }

// SetNamespace implements the metav1.Object interface.
func (e *Endpoints) SetNamespace(namespace string) {}

// GetName implements the metav1.Object interface.
func (e *Endpoints) GetName() string { return e.Name }

// SetName implements the metav1.Object interface.
func (e *Endpoints) SetName(name string) {}

// GetResourceVersion implements the metav1.Object interface.
func (e *Endpoints) GetResourceVersion() string { return e.Version }

// SetResourceVersion implements the metav1.Object interface.
func (e *Endpoints) SetResourceVersion(version string) {}
