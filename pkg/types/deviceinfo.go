package types

import (
	"log"
	"sync"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type DeviceInfo struct {
	idx    int
	podMap map[types.UID]*v1.Pod
	rwmu   *sync.RWMutex
	isUsed bool
}

func (d *DeviceInfo) IsUsed() bool {
	return d.isUsed
}

func (d *DeviceInfo) GetDevId() int {
	return d.idx
}

func newDeviceInfo(index int) *DeviceInfo {
	return &DeviceInfo{
		idx:    index,
		podMap: map[types.UID]*v1.Pod{},
		rwmu:   new(sync.RWMutex),
	}
}

func (d *DeviceInfo) AddPod(pod *v1.Pod) {
	log.Printf("debug: dev.addPod() Pod %s in ns %s with the GPU ID %d will be added to device map",
		pod.Name,
		pod.Namespace,
		d.idx)
	d.rwmu.Lock()
	defer d.rwmu.Unlock()
	d.podMap[pod.UID] = pod
	d.isUsed = true
	log.Printf("debug: dev.addPod() is success", )
}

func (d *DeviceInfo) RemovePod(pod *v1.Pod) {
	log.Printf("debug: dev.removePod() Pod %s in ns %s with the GPU ID %d will be removed from device map",
		pod.Name,
		pod.Namespace,
		d.idx)
	d.rwmu.Lock()
	defer d.rwmu.Unlock()
	delete(d.podMap, pod.UID)
	d.isUsed = false
	log.Printf("debug: dev.removePod() us success")
}
