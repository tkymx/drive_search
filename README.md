# drive_search
search target file identifier and full path  from google drive 

prepare 
```
go get -u google.golang.org/api/drive/v3
go get -u golang.org/x/oauth2/google
```

note : put credentials.json in same directory with drive_search
refer to https://developers.google.com/drive/api/v3/quickstart/go?hl=ja

running
```
go run drive_search.go -n SearchWord

or

./drive_search -n SearchWord
```

