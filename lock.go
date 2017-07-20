package ibclient

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	timeout     int32  = 60 // in seconds
	freeLockVal string = "Available"
)

type Lock interface {
	Lock() bool
	UnLock(force bool) bool
}

type NVLocker struct {
	name          string
	objMgr        *ObjectManager
	LockEA        string
	LockTimeoutEA string
}

func (l *NVLocker) createLockRequest() *MultiRequest {

	req := NewMultiRequest(
		[]*RequestBody{
			&RequestBody{
				Method: "GET",
				Object: "networkview",
				Data: map[string]interface{}{
					"name":         l.name,
					"*" + l.LockEA: freeLockVal,
				},
				Args: map[string]string{
					"_return_fields": "extattrs",
				},
				AssignState: map[string]string{
					"NET_VIEW_REF": "_ref",
				},
				Discard: true,
			},
			&RequestBody{
				Method: "PUT",
				Object: "##STATE:NET_VIEW_REF:##",
				Data: map[string]interface{}{
					"extattrs+": map[string]interface{}{
						l.LockEA: map[string]string{
							"value": l.objMgr.TenantID,
						},
						l.LockTimeoutEA: map[string]int32{
							"value": int32(time.Now().Unix()),
						},
					},
				},
				EnableSubstitution: true,
				Discard:            true,
			},
			&RequestBody{
				Method: "GET",
				Object: "##STATE:NET_VIEW_REF:##",
				Args: map[string]string{
					"_return_fields": "extattrs",
				},
				AssignState: map[string]string{
					"DOCKER-ID": "*" + l.LockEA,
				},
				EnableSubstitution: true,
				Discard:            true,
			},
			&RequestBody{
				Method: "STATE:DISPLAY",
			},
		},
	)

	return req
}

func (l *NVLocker) createUnlockRequest(force bool) *MultiRequest {

	getData := map[string]interface{}{"name": l.name}
	if !force {
		getData["*"+l.LockEA] = l.objMgr.TenantID
	}

	req := NewMultiRequest(
		[]*RequestBody{
			&RequestBody{
				Method: "GET",
				Object: "networkview",
				Data:   getData,
				Args: map[string]string{
					"_return_fields": "extattrs",
				},
				AssignState: map[string]string{
					"NET_VIEW_REF": "_ref",
				},
				Discard: true,
			},
			&RequestBody{
				Method: "PUT",
				Object: "##STATE:NET_VIEW_REF:##",
				Data: map[string]interface{}{
					"extattrs+": map[string]interface{}{
						l.LockEA: map[string]string{
							"value": freeLockVal,
						},
					},
				},
				EnableSubstitution: true,
				Discard:            true,
			},
			&RequestBody{
				Method: "PUT",
				Object: "##STATE:NET_VIEW_REF:##",
				Data: map[string]interface{}{
					"extattrs-": map[string]interface{}{
						l.LockTimeoutEA: map[string]interface{}{},
					},
				},
				EnableSubstitution: true,
				Discard:            true,
			},
			&RequestBody{
				Method: "GET",
				Object: "##STATE:NET_VIEW_REF:##",
				Args: map[string]string{
					"_return_fields": "extattrs",
				},
				AssignState: map[string]string{
					"DOCKER-ID": "*" + l.LockEA,
				},
				EnableSubstitution: true,
				Discard:            true,
			},
			&RequestBody{
				Method: "STATE:DISPLAY",
			},
		},
	)

	return req
}

func (l *NVLocker) Lock() bool {
	logrus.Debugf("Creating lock on network niew %s\n", l.name)
	req := l.createLockRequest()
	res, err := l.objMgr.CreateMultiObject(req)

	if err != nil {
		logrus.Debugf("Failed to create lock on network view %s: %s\n", l.name, err)

		//Check for Lock Timeout
		nw, err := l.objMgr.GetNetworkView(l.name)
		if err != nil {
			logrus.Debugf("Failed to get the network view object for %s : %s\n", l.name, err)
			return false
		}

		if t, ok := nw.Ea[l.LockTimeoutEA]; ok {
			if int32(time.Now().Unix())-int32(t.(int)) > timeout {
				logrus.Debugln("Lock is timed out. Forcefully acquiring it.")
				//remove the lock forcefully and acquire it
				l.UnLock(true)
				// try to get lock again
				return l.Lock()
			}
		}
		return false
	}

	dockerID := res[0]["DOCKER-ID"]
	if dockerID == l.objMgr.TenantID {
		logrus.Debugln("Got the lock !!!")
		return true
	}

	return false
}

func (l *NVLocker) UnLock(force bool) bool {
	// To unlock set the Docker-Plugin-Lock EA of network view to Available and
	// remove the Docker-Plugin-Lock-Time EA
	req := l.createUnlockRequest(force)
	res, err := l.objMgr.CreateMultiObject(req)

	if err != nil {
		logrus.Debugf("Failed to release lock on Network View %s: %s\n", l.name, err)
		return false
	}

	dockerID := res[0]["DOCKER-ID"]
	if dockerID == freeLockVal {
		logrus.Debugln("Removed the lock!")
		return true
	}

	return false
}

func GetNVLock(netViewName string, objMgr *ObjectManager, lockEA string, lockTimeoutEA string) (Lock, error) {

	// verify if network view exists and has EA for the lock
	nw, err := objMgr.GetNetworkView(netViewName)
	if err != nil {
		msg := fmt.Sprintf("Failed to get the network view object for %s : %s\n", netViewName, err)
		logrus.Debugf(msg)
		return nil, fmt.Errorf(msg)
	}

	if _, ok := nw.Ea[lockEA]; !ok {
		err = objMgr.UpdateNetworkViewEA(nw.Ref, EA{lockEA: freeLockVal}, nil)
		if err != nil {
			return nil, fmt.Errorf("Failed to Update Network view with Lock EA")
		}
	}

	l := &NVLocker{name: netViewName, objMgr: objMgr, LockEA: lockEA, LockTimeoutEA: lockTimeoutEA}
	retryCount := 0
	for {
		// Get lock on the network view
		lock := l.Lock()
		if lock == true {
			// Got the lock.
			logrus.Debugf("Got the lock on Network View %s\n", netViewName)
			return l, nil
		} else {
			// Lock is held by some other agent. Wait for some time and retry it again
			if retryCount >= 10 {
				return nil, fmt.Errorf("Failed to get Lock on Network View %s", netViewName)
			}

			retryCount++
			logrus.Debugf("Lock on Network View %s not free. Retrying again %d out of 10.\n", netViewName, retryCount)
			// sleep for random time (between 1 - 10 seconds) to reduce collisions
			time.Sleep(time.Duration(rand.Intn(9)+1) * time.Second)
			continue
		}
	}
}
