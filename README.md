`ijq` stands for "interactive jq", the super-powered [jq](https://stedolan.github.io/jq/) REPL with support for:

  * automatic variable assignment -- so you won't lose your history!
  * global function definition -- define functions now, use them later in the same session!
  * `import` and `include` statements!

Example session
---------------

```json
fiatjaf@cantillon ~> echo '{"numbers": [1,2]}' > data.json
fiatjaf@cantillon ~> echo 'def sum(a;b): a + b;' > math.jq
fiatjaf@cantillon ~> ijq data.json 
(./jq)| import "math" as math
(./jq)| math::sum(.numbers[0]; .numbers[1])
3 as $v1
(./jq)| $v0 | .sum = $v1
{
  "numbers": [
    1,
    2
  ],
  "sum": 3
} as $v2
(./jq)| def addtag(tagname): .tags = (.tags // []) | .tags += [tagname]
(./jq)| $v2 | addtag("silly-math")
{
  "numbers": [
    1,
    2
  ],
  "sum": 3,
  "tags": [
    "silly-math"
  ]
} as $v3
(./jq)| addtag("trivial")
{
  "numbers": [
    1,
    2
  ],
  "sum": 3,
  "tags": [
    "silly-math",
    "trivial"
  ]
} as $v4
(./jq)| ⏎
```

Installation
------------

Install using [`go install`](https://go.dev/ref/mod#go-install):
```
go install github.com/fiatjaf/ijq/...@latest
```

Recommended:

```
sudo apt-get install rlwrap # or whatever, but please install rlwrap
```

Then

```
ijq [file]
```

FAQ
---

### My commands are failing and I don't understand why!

Use the special `debug` command, you'll get the full filter that is being passed to `jq` and will be able to know what is happening. If it is a bug on `ijq` report it [here](https://github.com/fiatjaf/ijq/issues) please!

---

[![Mentioned in Awesome jq](https://awesome.re/mentioned-badge.svg)](https://github.com/fiatjaf/awesome-jq)
