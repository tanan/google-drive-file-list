# google-drive-file-list
List files with a given folder id in Google Drive. This code is written by Go and use Google Drive API.

## prerequisites

We need to add Drive scope when using Google Drive API. Please run the following command:

```
gcloud auth application-default login --scopes "https://www.googleapis.com/auth/drive.readonly,https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/cloud-platform,openid"
```

## usage
- your drive

```
go run main.go node.go -f [folderID]
```

- shared drive

```
go run main.go node.go -f [folderID] -d [driveID]
```

## how to get folder ID and dirve ID

For example, if the link is "https://drive.google.com/drive/folders/ABC123456789", then the folder ID is "ABC123456789".
Drive ID can get like above as a root folder.
