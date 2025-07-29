# for upload providers these are the functions we need

- a function given a multipart.header will upload the file and give us back the url{Or an update dto Containing url}.
- a function given a name of the file, it will delete the file

```go
    type UploadDto struct {
        Name string // the unique name fo the file
        FileType string // the filetype got from the header
        Url string  // the url which we will get this file, should have ip or domain
        Size int64 //the actual size of the file
        Hash string // the unique hash of the file
        UserId string // the id of the user who uploaded the file
        IsMainImg bool
    }

```

- we could have functions on files like

- Users file calculation query
  - the total number of files a user has
  - the total size of files a user has
  - the count of files a users has per filetype

## case of different model file Uploads

- we let usually users to upload files for usually 3 models
  - items(ecom), profile photo, category or there models with only one image
- so we can make separate endpoints for each of those, who have access to the database with its own image upload rules
- for admins we could make a separate upload endpoints.

## case of Uploads with cover image

- we could
