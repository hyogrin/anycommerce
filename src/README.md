## Application codes for Anycommerce

for each directories (e.g. carts):
```
rm -rf .git
git init
git remote add origin codecommit::ap-northeast-2://${PWD##*/}
git add . 
git commit -m "Initial commit." 
git branch -M main
git push -u origin main
```