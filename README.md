# Go 程式碼基本架構

## 手刻方法

### 直接從 Github 下載模板
```bash=
git clone https://github.com/POABOB/go-base-template
```

## IMOOC 購買課程方法
### 1.快速創建程式碼倉庫請使用下方命令
```bash=
sudo docker run --rm -v $(pwd): $(pwd) -w  $(pwd) -e ICODE=xxxxxx cap1573/cap-tool new git.imooc.com/cap1573/base

注意：
1.sudo 如果是 Mac 系統提醒輸入的密碼是本機的密碼。
2.以上命令ICODE=xxxxxx 中 "xxxxxx" 為個人購買的 icode 碼。
3.icode 馬在購買完課程後，請使用電腦點擊進入學習課程頁面。
4.請勿多人使用同一個 icode 碼（會被慕課網封鎖）。
5.這裡 git.imooc.com/cap1573/base 倉庫 名字需要和 go mod 一致
```
 

### 2.根據 proto 自動生成 go 基礎程式碼
```bash=
make proto
```

### 3.根據程式碼編譯現有的 Go 程式  
```bash=
make build
```
代码执行后会产生 base 二进制文件
程式碼執行完成後會產生 base 二進位制文件

### 4.編譯執行二進位制文件
```bash=
make docker
```
編譯成功後會自動產生 base:latest 鏡像
可使用 docker images | grep base 查看是否產生

### 5.本課程使用go-micro v3 版本作為微服務開發框架
框架地址：https://github.com/asim/go-micro
