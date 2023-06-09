/*
 * Copyright (c) The One True Way 2023. Apache License 2.0. The authors accept no liability, 0 nada for the use of this software.  It is offered "As IS"  Have fun with it!!
 */

package metrics

import "testing"

func TestServerFactor(t *testing.T) {
	InitProxyServerMetrics()
	IncTotalFailedRequests("401")
	IncTotalRequests()
}
func TestRelayletFactor(t *testing.T) {
	InitProxyletMetrics()
	IncTotalFailedRequests("401")
	IncTotalRequests()
}
