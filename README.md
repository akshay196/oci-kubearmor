# OCI KubeArmor

```
$ go run main.go block-k8s-sa-token.yaml localhost:5001/block-sa:latest
file descriptor for block-k8s-sa-token.yaml: {application/vnd.cncf.kubearmor.policy.layer.v1.yaml sha256:fe987b7f35b13bb9f5d71a3bb94e4e4f4606cee950a3fe9e61a3187b2198cbe7 542 [] map[org.opencontainers.image.title:block-k8s-sa-token.yaml] [] <nil> }
manifest descriptor: {application/vnd.oci.image.manifest.v1+json sha256:07c4669eecdbbce1708d96f3026e4c3b19c2b9e16c8762119d5662c6015bb43f 561 [] map[org.opencontainers.image.created:2023-06-11T05:38:55Z] [] <nil> application/vnd.cncf.kubearmor.config.v1+json}
```


```
$ crane manifest localhost:5001/block-sa:latest | jq
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.cncf.kubearmor.config.v1+json",
    "digest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
    "size": 2
  },
  "layers": [
    {
      "mediaType": "application/vnd.cncf.kubearmor.policy.layer.v1.yaml",
      "digest": "sha256:fe987b7f35b13bb9f5d71a3bb94e4e4f4606cee950a3fe9e61a3187b2198cbe7",
      "size": 542,
      "annotations": {
        "org.opencontainers.image.title": "block-k8s-sa-token.yaml"
      }
    }
  ],
  "annotations": {
    "org.opencontainers.image.created": "2023-06-11T05:38:55Z"
  }
}
```

```
$ crane blob localhost:5001/block-sa:latest@sha256:fe987b7f35b13bb9f5d71a3bb94e4e4f4606cee950a3fe9e61a3187b2198cbe7
apiVersion: security.kubearmor.com/v1
kind: KubeArmorPolicy
metadata:
  name: ksp-wordpress-block-sa
  namespace: wordpress-mysql
spec:
  severity: 7
  selector:
    matchLabels:
      app: wordpress
  file:
    matchDirectories:
    - dir: /run/secrets/kubernetes.io/serviceaccount/
      recursive: true

      # cat /run/secrets/kubernetes.io/serviceaccount/token
      # curl https://$KUBERNETES_PORT_443_TCP_ADDR/api --insecure --header "Authorization: Bearer $(cat /run/secrets/kubernetes.io/serviceaccount/token)"

  action:
    Block
```
