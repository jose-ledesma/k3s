package cloudprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rancher/k3s/pkg/version"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

var (
	InternalIPAnnotation = version.Program + ".io/internal-ip"
	ExternalIPAnnotation = version.Program + ".io/external-ip"
	HostnameAnnotation   = version.Program + ".io/hostname"
)

func (k *k3s) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return cloudprovider.NotImplemented
}

func (k *k3s) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (k *k3s) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	return true, nil
}

func (k *k3s) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	if k.nodeInformerHasSynced == nil || !k.nodeInformerHasSynced() {
		return "", errors.New("Node informer has not synced yet")
	}

	_, err := k.nodeInformer.Lister().Get(string(nodeName))
	if err != nil {
		return "", fmt.Errorf("Failed to find node %s: %v", nodeName, err)
	}
	return string(nodeName), nil
}

func (k *k3s) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return true, cloudprovider.NotImplemented
}

func (k *k3s) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	_, err := k.InstanceID(ctx, name)
	if err != nil {
		return "", err
	}
	return version.Program, nil
}

func (k *k3s) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	return "", cloudprovider.NotImplemented
}

func (k *k3s) NodeAddresses(ctx context.Context, name types.NodeName) ([]corev1.NodeAddress, error) {
	addresses := []corev1.NodeAddress{}
	if k.nodeInformerHasSynced == nil || !k.nodeInformerHasSynced() {
		return nil, errors.New("Node informer has not synced yet")
	}

	node, err := k.nodeInformer.Lister().Get(string(name))
	if err != nil {
		return nil, fmt.Errorf("Failed to find node %s: %v", name, err)
	}
	// check internal address
	if node.Annotations[InternalIPAnnotation] != "" {
		for _, address := range strings.Split(node.Annotations[InternalIPAnnotation], ",") {
			addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: address})
		}
	} else {
		logrus.Infof("Couldn't find node internal ip label on node %s", name)
	}

	// check external address
	if node.Annotations[ExternalIPAnnotation] != "" {
		for _, address := range strings.Split(node.Annotations[ExternalIPAnnotation], ",") {
			addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalIP, Address: address})
		}
	}

	// check hostname
	if node.Annotations[HostnameAnnotation] != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeHostName, Address: node.Annotations[HostnameAnnotation]})
	} else {
		logrus.Infof("Couldn't find node hostname label on node %s", name)
	}

	return addresses, nil
}

func (k *k3s) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]corev1.NodeAddress, error) {
	return nil, cloudprovider.NotImplemented
}
