# Deploying Fern Ecosystem with KubeVela

This guide walks you through deploying the complete Fern ecosystem (PostgreSQL database, fern-reporter, and fern-mycelium) using KubeVela on a k3d cluster.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [k3d](https://k3d.io/v5.6.0/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Helm](https://helm.sh/docs/intro/install/)
- [KubeVela CLI](https://kubevela.io/docs/installation/standalone)

## Setup

### 1. Create k3d cluster

```bash
k3d cluster create my-k3d-cluster --port "8080:8080@loadbalancer" --agents 3
```

### 2. Install KubeVela

```bash
helm repo add kubevela https://kubevela.github.io/charts
helm repo update
helm install --create-namespace -n vela-system kubevela kubevela/vela-core --wait
```

### 3. Install Cloud Native PostgreSQL

```bash
kubectl create namespace cnpg-system
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm repo update
helm install cnpg cnpg/cloudnative-pg -n cnpg-system --create-namespace --wait
```

### 4. Install Custom ComponentDefinitions

Install the KubeVela CLI:
```bash
curl -fsSl https://kubevela.io/script/install.sh | bash
```

Apply the custom component definitions:
```bash
vela def apply cnpg.cue
vela def apply gateway.cue
```

## Deployment

### 1. Create namespace
```bash
kubectl create namespace fern
```

### 2. Deploy the application
```bash
kubectl apply -f ./docs/kubevela/vela.yaml
```

### 3. Verify deployment
```bash
kubectl get all -n fern
kubectl get application -n fern
vela status fern-ecosystem -n fern
```

## Access Applications

The applications will be available at:
- **fern-reporter**: http://fern-reporter.local
- **fern-mycelium**: http://fern-mycelium.local

To access locally, add these entries to your `/etc/hosts`:
```
127.0.0.1 fern-reporter.local
127.0.0.1 fern-mycelium.local
```

## Configuration

### Database Configuration
The PostgreSQL database is configured with:
- 1 instance
- 1Gi storage
- Database name: `fern`
- Credentials automatically generated and shared via secrets

### Environment Variables
Both applications receive the following database credentials:
- `FERN_USERNAME`: Database username
- `FERN_PASSWORD`: Database password  
- `FERN_HOST`: Database host
- `FERN_PORT`: Database port
- `FERN_DATABASE`: Database name
- `DB_URL`: Complete database connection URI (for fern-mycelium)

## Troubleshooting

### Check application status
```bash
vela status fern-ecosystem -n fern
kubectl describe application fern-ecosystem -n fern
```

### Check component logs
```bash
kubectl logs -n fern deployment/fern-reporter
kubectl logs -n fern deployment/fern-mycelium
```

### Check database status
```bash
kubectl get cluster -n fern
kubectl describe cluster postgres -n fern
```

## Cleanup

```bash
kubectl delete application fern-ecosystem -n fern
kubectl delete namespace fern
k3d cluster delete my-k3d-cluster
```

## Contributing

We welcome contributions and customizations to this deployment configuration. Please feel free to submit issues or pull requests with improvements.