# Helm Charts for AnyCommerce

Applications ready to launch on Kubernetes using [Helm](https://github.com/helm/helm).

## Using Chart

- Install a chart: `helm install my-release ./anycommerce/`
- Upgrade your application: `helm upgrade my-release ./anycommerce/`
- Put your charts under the anycommerce directory.

## Example

Update `charts/anycommerce/values.yaml`
```
# (Optional) Change these to your own repository url.
carts:
  image:
    repository: 269550163595.dkr.ecr.ap-northeast-2.amazonaws.com/anycommerce/carts

order:
  image:
    repository: 269550163595.dkr.ecr.ap-northeast-2.amazonaws.com/anycommerce/order

product:
  image:
    repository: 269550163595.dkr.ecr.ap-northeast-2.amazonaws.com/anycommerce/product

recommendation:
  image:
    repository: 269550163595.dkr.ecr.ap-northeast-2.amazonaws.com/anycommerce/recommendation

user:
  image:
    repository: 269550163595.dkr.ecr.ap-northeast-2.amazonaws.com/anycommerce/user

# (Required) Change this to your own ALB arn. You can find in ACM dashboard.
    alb.ingress.kubernetes.io/certificate-arn: >-
      arn:aws:acm:ap-northeast-2:994385297847:certificate/cb9c46f6-9737-4c87-802c-341c26c7b3df
```


```bash
cd charts/anycommerce

# Install Helm chart for the environment.
kubectl create namespace anycommerce
helm install anycommerce . -n anycommerce
```