package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/skratchdot/open-golang/open"
)

const (
	RC_FILE             = ".chickenheadrc"
	DEFAULT_SNIPPET_DIR = ".chickenhead"
)

var (
	app  *cli.App
	conf *Config
)

type Config struct {
	SnippetDirectory string
	Editor           string
}

func (c *Config) String() string {
	return fmt.Sprintf("%#v\n", c)
}

func setSubCommands() {
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a new snippet",
			Action:  add,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "s",
					Usage: "read contents from Standard Input",
				},
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete a snippet",
			Action:  delete,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "f",
					Usage: "force delete",
				},
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "edit a snippet",
			Action:  edit,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list up available snippets",
			Action:  list,
		},
		{
			Name:    "view",
			Aliases: []string{"v"},
			Usage:   "view a snippet",
			Action:  view,
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search snippets",
			Action:  search,
		},
	}
}

// add is a subcommand to create a new snippet.
func add(c *cli.Context) {
	snippetName := c.Args().First()
	if len(snippetName) == 0 {
		log.Fatal("please enter snippet name")
	}

	s := NewSnippet(conf.SnippetDirectory, snippetName)

	if s.Exists() {
		log.Fatal("This snippet already exists.")
	}

	stdinFlag := c.Bool("s")
	switch {
	case stdinFlag:
		// read contents from standard input
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		// write contents to snippet
		s.Write(buf)
	default:
		f, err := s.Create()
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		openSnippetWithEditor(s)
	}
}

// delete is a subcommand to delete a snippet.
func delete(c *cli.Context) {
	force := c.Bool("f")
	snippetName := c.Args().First()
	if len(snippetName) == 0 {
		log.Fatal("please enter snippet name")
	}

	s := NewSnippet(conf.SnippetDirectory, snippetName)

	if !s.Exists() {
		log.Fatal("This snippet doesn't exists.")
	}

	if !force {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("sure you want to delete %s [yn]?")
		ret, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if ret != "y\n" {
			return
		}
	}

	if err := s.Delete(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s has been deleted.\n", snippetName)
}

// edit is subcommand to edit a snippet with an editor.
func edit(c *cli.Context) {
	snippetName := c.Args().First()
	if len(snippetName) == 0 {
		log.Fatal("please enter snippet name")
	}

	s := NewSnippet(conf.SnippetDirectory, snippetName)
	if !s.Exists() {
		log.Fatal("snippet not found :", s.Path)
	}

	openSnippetWithEditor(s)
}

// list is a subcommand to list up snippets.
// It just finds snippet files in the snippet directory and listed them.
func list(c *cli.Context) {
	var pattern *regexp.Regexp
	var err error
	query := c.Args().First()
	if len(query) > 0 {
		pattern, err = regexp.Compile(fmt.Sprintf(".*%s.*", query))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = filepath.Walk(
		conf.SnippetDirectory,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(conf.SnippetDirectory, path)

			if pattern != nil {
				if pattern.MatchString(rel) {
					fmt.Println(rel)
				}
				return nil
			}

			fmt.Println(rel)
			return nil
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}

// view is a subcommand to show contents of a snippet.
func view(c *cli.Context) {
	snippetName := c.Args().First()
	if len(snippetName) == 0 {
		log.Fatal("please enter snippet name")
	}

	s := NewSnippet(conf.SnippetDirectory, snippetName)
	if !s.Exists() {
		log.Fatal("snippet not found :", s.Path)
	}

	contents, err := s.ReadContents()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(contents)
}

// serch is a subcommand to searche snippets whose contents contains the query with pt or ag.
func search(c *cli.Context) {
	var cmd *exec.Cmd
	query := c.Args().First()

	if len(query) == 0 {
		log.Fatal("please enter search query")
	}

	switch {
	case commandExists("pt"):
		cmd = exec.Command("pt", "-i", query, conf.SnippetDirectory)
	case commandExists("ag"):
		cmd = exec.Command("ag", query, conf.SnippetDirectory)
	default:
		log.Fatal("the_platinum_searcher(pt) or the_silver_searcher(ag) is required.")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// When query is not matched, both commands return an error,
	// so the error is ignored here.
	_ = cmd.Run()
}

func loadConfig(homeDir string) *Config {
	rcFilePath := filepath.Join(homeDir, RC_FILE)
	cf := new(Config)

	// read config values if a rc file exists.
	if exists(rcFilePath) {
		f, err := os.Open(rcFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := toml.Decode(string(b), cf); err != nil {
			log.Fatal(err)
		}
	}

	// set default values if empty.
	if len(cf.SnippetDirectory) == 0 {
		cf.SnippetDirectory = filepath.Join(homeDir, DEFAULT_SNIPPET_DIR)
	}

	return cf
}

func openSnippetWithEditor(s *Snippet) {
	if len(conf.Editor) > 0 {
		err := open.RunWith(s.Path, conf.Editor)
		if err != nil {
			log.Fatalf("can't open %s with %s. please check your settings.", s.Path, conf.Editor)
		}

	} else {
		err := open.Run(s.Path)
		if err != nil {
			log.Fatalf("can't open %s. please check your settings.", s.Path)
		}
	}
}

func commandExists(cmdName string) bool {
	_, err := exec.LookPath(cmdName)
	return err == nil
}

// exists checks a file exists or not.
func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func mkDir(dirPath string) {
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
}

func getHomeDir() string {
	var err error
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir
}

func main() {
	conf = loadConfig(getHomeDir())
	app = cli.NewApp()
	app.Name = "chickenhead"
	app.Usage = "simple CLI snippet tool."
	setSubCommands()
	mkDir(conf.SnippetDirectory)
	app.Run(os.Args)
}
