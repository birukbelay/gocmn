# some usefull commands

- tagging:

> git tag v1.0
> git tag -a v1.0 -m "Release version 1.0"

- View Tags

> git tag

- To list tags matching a pattern:

> git tag -l "v1.*"

- Tag a Specific Commit

> git tag -a v1.0 <commit-hash> -m "Tag message"

- Push a single tag:

> git push origin v1.0

- Push all tags:

> git push origin --tags

Delete a Tag

- Locally:

> git tag -d v1.0

- On remote:

> git push origin --delete tag v1.0
