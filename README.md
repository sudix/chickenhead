# chickenhead

Description
=============

`chickenhead` is a simple CLI snippet tool.

In japanese, `とりあたま`(chicken head) means having a bad memory.
And I have a bad memory.

Snippet Files
=============

Snippets are stored as text files under group directories.
Groups can be specified with "/".

When a snippet is created,
group directories are created under SNIPPET_DIRECTORY(e.g. /home/sudix/.chickenhead),
and then a snippet file is created in the bottom of the group directories.

This command

```
$ chickenhead add -s how/to/buy/books
Go Amazon.com!
```

create a snippet like below.

```
SNIPPET_DIRECTORY
├── how                  # group directory 1
│   └── to              # group directory 2
│       └── buy         # group directory 3
│           └── books   # snippet file (that contens is "Go Amazon.com!")
```

Configuration
=============

`chickenhead` has two configuration values.

* `SnippetDirectory` - A directory to store snippet files. Default value is `$HOME/.chickenhead`.
* `Editor` - An editor that you want to open or edit snippet files.

You can specify these values in a config file.
That config file must be `$HOME/.chickenheadrc`.

### Sample Configuration

` .chickenheadrc`

```
SnippetDirectory = "/Users/sudix/Dropbox/chickenhead"
Editor = "CotEditor.app"
```

INSTALLATION
=============

```
$ go get github.com/sudix/chickenhead
```

COMMANDS
=============

### add

Add a new snippet.

Without a flag, new snippet file is created and open it with an editor.

```
$ chickenhead add go/http/static_server
```

With `s` flag, read contents from standard input and write it to new snippet file.

```
$ chickenhead add -s go/http/static_server
foo
bar
baz
```

### delete

Delete a specified snippet.

```
$ chickenhead delete go/http/static_server
```

When you want to delete a snippet without confirmatin, "-f" (force delete flag) is convenient.

```
$ chickenhead delete -f go/http/static_server
```

### edit

Edit a specified snippet.

```
$ chickenhead edit go/http/static_server
# This opens default editor with given name.
```

### list

List up available snippets

```
$ chickenhead list
go/http_server
go/http_client
go/static_server
go/open_file
ruby/how_to_open_file
sed/git_origin
        .
        .
        .
```

If a query arguments is given, only snippets whose name contain the query are listed.

```
$ chickenhead list go
go/http/static_server
go/http/client
go/io/read_file
```

### view

```
$ chickenhead view go/http/static_server
```

### search

Search snippets whose contents includes the query.
This command use [the_platinum_searcher](https://github.com/monochromegane/the_platinum_searcher) or [the_silver_searcher](https://github.com/ggreer/the_silver_searcher)  internally,
so you need to install one of them.

```
$ chickenhead serach http
```

Aliases
=============

`chickenhead` is too long to remember, so I recommend to use aliases.
Using list command with [peco](https://github.com/peco/peco) or [percol](https://github.com/mooz/percol) is very useful.

### Alias examples

```
alias cha="chickenhead add"
alias chd="chickenhead delete"
alias chl="chickenhead list | peco"
alias chv="chickenhead list | peco | xargs chickenhead view"
alias che="chickenhead list | peco | xargs chickenhead edit"
alias chs="chickenhead search"
```
