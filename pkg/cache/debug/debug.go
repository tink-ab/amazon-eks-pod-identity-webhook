package debug

import (
	"fmt"
	"net"
	"net/http"

	"github.com/aws/amazon-eks-pod-identity-webhook/pkg/cache"
	"k8s.io/klog"
)

type Dumper struct {
	Cache cache.ServiceAccountCache
}

func (c *Dumper) Handle(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || !net.ParseIP(host).IsLoopback() {
		http.Error(w, "Forbidden: debug endpoints are only accessible from localhost", http.StatusForbidden)
		return
	}

	res := c.Cache.ToJSON()
	if _, err := w.Write([]byte(res)); err != nil {
		klog.Errorf("Can't dump cache contents: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
