# cursor based paginataion

## going forward

lets say we have records
with values (num, id) and id must always be unique and num can be repeated

### sort asc

1a, 2b, 3a, 3b, 3c, 3d, 4a, 5a

if we want cursor based pagination let say limit is 3

so we sort asc,

so we will get, 1a, 2b, 3a

"(%s > ?) OR (%s = ? AND id > ?)
"(num > 3) OR (num = 3 AND id > a)

so we get: 3b, 3c, 3d

our next cursor will be, 3d
and our prev cursor will be 3b

### sort desc

- reminder: 3a, 3b, 3c are always coming in the same order, because the id is always asc
- ?? will this help us in only creating 1 index, id ascending
5a, 4a, 3e,| 3d, 3c, 3b |, 3a, 2b, 2c, 1a

we want the next three after 3e: 3d, 3c, 3b

- "(%s < ?) OR (%s = ? AND id < ?)"

- "(num < 3) OR (num = 3 AND id < e)"

---

## if we want to go backward

### asc

----  actual query sort: `num asc, id asc`

1a, 2b, 3a,| 3b, 3c, 3d,| `3e`, 4a, 5a

we are looking for values before 3e:

num < 3e, but if we are sorted asc, it will ruin the pagination
so we need to sort it desc

& in order to keep the order the id needs to be opposite
we want(3d, 3c, 3b) reversed as (3b, 3c, 3d) in that order so

----------- query reverse sort `num desc, id desc`

5a, 4a, `3e`, | 3d, 3c, 3b, |3a, 2b , 1a

"(num < 3) OR (num = 3 AND id < e)"
"(%s < ?) OR (%s = ? AND id < ?)"
we get 3d, 3c, 3b -> reversed it will be: 3b, 3c, 3d

our fwd cursor = 3d
our bwd cursor = 3b

### desc

Orignal Sort: `num desc, id desc`

5a, 4a, 3e, | 3d, 3c, 3b, |`3a`, 2b , 1a

we will need values before `3a`; which are-- 3d, 3c, 3b
we will want it to be reversed, sort by: num asc, id acs

1a, 2b, `3a`,| 3b, 3c, 3d,| 3e, 4a, 5a

"(num > 3) OR (num = 3 AND id > a)"
"(%s > ?) OR (%s = ? AND id > ?)"

nextCursor= 3b:
prevCursor= 3d

--- when id is always asc

- Sort: num desc, id asc
   5a, 4a, 3a, | 3b, 3c, 3d, |3e, 2b , 1a
