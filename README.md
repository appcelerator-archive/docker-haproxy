###amp-haproxy-controller Prototype


It's a HAProxy image containing a controller which is able to update HAProxy configuration when some etcd keys are updated

# tags

- latest

# Setting

environment variable: 'ETCD_ENDPOINTS' should be equal to the ETCD endpoints list, format: "host1:port1, host2:port2, ..."
