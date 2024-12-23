# Alertfy - Bridging AlertManager to NTFY

## Getting Started

```
git clone https://github.com/murtaza-u/alertfy
cd helm
```

Edit `values.yaml`

The file `values.yaml` contains everything you can configure. Once you have
modified the values to your preferences, run the `helm install` command below
to deploy Alertfy in your cluster:

```
NAMESPACE="monitoring"

helm upgrade alertfy . \
    --install \
    --namespace "$NAMESPACE" \
    --create-namespace \
    --values values.yaml
```
