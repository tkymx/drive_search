# drive_search
search target file identifier and full path  from google drive 

# drive_download
download target file from google drive 

prepare 
```
go get -u google.golang.org/api/drive/v3
go get -u golang.org/x/oauth2/google
```

note : put credentials.json in same directory with drive_search
refer to https://developers.google.com/drive/api/v3/quickstart/go?hl=ja

running 
* in each directory
```
go run drive_search.go -n SearchWord
go run drive_download.go -i 0BzWl_pRfWl8wRUQ5Ds52c2ctRWc -b ../

or

sh build.sh 
./drive_search -n SearchWord
./drive_download.go -i 0BzWl_pRfWl8wRUQ5Ds52c2ctRWc -b ../
```




