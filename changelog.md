# change Log

> git tag     # to list the tags
> git tag -a v1.0 -m "Release version 1.0"
> git push origin --tags
> git push origin v1.0

## 1.1

- major rewrite

### 1.1.2

- removed `sel = append(sel, options.Preloads...)` from generic.query-helpers.DebugPreloadSelect
- changed from . to - for file separtaion, to use `.` for naming functions inside files

### 1.1.3

- [x] add getKeys and getValues to the util/set fucntions
- [x] add clause.Returning{} for the add function
- [x] add the error to the huma Greturn funciton, which takes ...err params

### 1.1.6

- added add multiple vals to set and remove multiple vals from set
- added more colors and functions to logger.
- added firebase service with push notification service
