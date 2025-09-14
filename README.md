# EchoConfig Operator

This is a simple Kubernetes Operator built with [Kubebuilder](https://book.kubebuilder.io/).  
It demonstrates how to extend the Kubernetes API with a **Custom Resource Definition (CRD)** and a **Controller**.

---

## 🚀 What it Does

- Defines a custom resource **EchoConfig** (`echoconfigs.demo.deet.dev/v1alpha1`).
- For each `EchoConfig`, the controller:
  - Creates/updates a `ConfigMap` named `echo-<name>` in the same namespace.
  - Writes the CR’s `spec.message` into `data.message` of the ConfigMap.
- Updates CR `status` with the name of the managed ConfigMap.

👉 Example:  
```yaml
apiVersion: demo.deet.dev/v1alpha1
kind: EchoConfig
metadata:
  name: hello
spec:
  message: "Hello from Deet!"
```

Produces this ConfigMap automatically:  
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: echo-hello
  namespace: default
data:
  message: "Hello from Deet!"
```

---

## ⚙️ How it Works

- **CRD (Custom Resource Definition):** Declares a new Kubernetes resource type (`EchoConfig`) with `spec.message` and `status.configMapName`.  
- **Controller:** Watches for `EchoConfig` objects and reconciles the desired state by creating/updating a ConfigMap.  

Analogy:  
> A CRD without a controller is like a **boat without someone rowing it** , it just floats.  
> The Controller is the rower who makes it move.  

---

## 🛠️ Prerequisites

- Go 1.21+  
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html)  
- Access to a Kubernetes cluster  
- Docker Hub account (for building/pushing images)  

---

## 📝 Development Workflow

### 1. Clone and Init
```bash
git clone https://github.com/cloudwithdeethesh/echoconfigmap.git
cd echoconfigmap
```

### 2. Install CRD into your cluster
```bash
make install
```

### 3. Run controller locally
```bash
make run
```

### 4. Apply a sample EchoConfig
```bash
kubectl apply -f config/samples/demo_v1alpha1_echoconfig.yaml
```

Check the ConfigMap:
```bash
kubectl get configmap echo-hello -o yaml
```

---

## 📦 Deploy as Kubernetes Deployment

To run this controller inside Kubernetes:

1. Build & push the image:
```bash
export IMG=docker.io/cloudwithdeethesh/echoconfigmap:v0.1.0
make docker-build IMG=$IMG
make docker-push IMG=$IMG
```

2. Deploy the controller:
```bash
make deploy IMG=$IMG
```

Now the operator runs in your cluster as a Deployment, watching for `EchoConfig` resources.

---

## 🔑 Takeaway

This project shows the **building blocks of Kubernetes Operators**:
- CRDs let you define new Kubernetes resource types.
- Controllers continuously reconcile those resources into real-world changes.
- Together, they extend Kubernetes with custom APIs.

---

## 📚 References

- [Kubebuilder Book](https://book.kubebuilder.io/)  
- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)  
