# SuperMux Roundtripper

SuperMux is a router plugin that handles HTTP and GRPC transport, it currently specifically targets Kubernetes.

## Use

```yaml
apiVersion: alpha.oscea.com/v1
Kind: Transpose
Metadata:
  name: myProxy
  namespace: default
spec:
  listener:
    name: mylistener
    package: github.com/oscea/transpose-plugins/listener/http
    spec: 
      port: 80
      ssl: false

  roundtrip:
    name: myroundtrip
    package: github.com/oscea/transpose-plugins/roundtrip/supermux
    spec:
      http:
      - path: "/"
        backend:
          serviceName: myservice
          servicePort: 80
```
