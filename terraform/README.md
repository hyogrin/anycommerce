## Terraform for AnyCommerce
**NOTE**
* Set `aws_account_id` and `domain_name` to yours in `variables.tf`.
* Register domain [here](https://supernova.amazon.dev/) if you are an AWS employee.
* Only two AZs are shown. We are deploying four AZs in this code.

### Create a Cluster
```bash
terraform init
terraform plan
terraform apply
aws eks update-kubeconfig --name anycommerce --region ap-northeast-2
kubectl get pods -A
```

### Deleting Cluster
```bash
terraform destroy
```