# mage-play

An evaluation of [Mage](https://magefile.org). The result is a project defining importable Mage targets that could be imported into other Magefiles, demonstrating two possible test approaches, 1) running a Magefile importing the targets, 2) unit testing the package.

## Targets

Mage targets are in `targets/`. They can be imported as per `test/mage.go`. Note that `magefile.go` exists to support
running with `go run` as used in the tests (`go run tests/mage.go ...`). This is to avoid the dependency on the `mage`
binary. If mage is installed on the system then `mage ...` in the `test/` directory would suffice.

### VersionBump

This target looks for the most recent tag on the current branch (HAED) by walking the commits. Bumps the version depending on the
commit message contents of HEAD, and tags HEAD. It then pushes the tag to `origin`. I'm on the fence as to whether it
should simply look for the latest of all tags, or only by walking the current branch.

I'm not convinced sing `go-git` is any better than simply using the `git` commands we're all familiar with, but it was
a fun learning experience.

## Why Mage?

Mage provides a similar end-user experience to Make, but written in Go. Make & Bash are default choice for build automation and has been battle tested throughout the years. They work well, and
are well understood by most engineers. Well, are they? I've been working with both for a number of years, and I'm not
ashamed to admit I stumble my way through Bash scripts and Makefiles with a chain of google searches.
Putting the time in to learn these tools properly would help, but I'm not convinced it would pay back the effort as I
just don't use them frequently enough.

With that said, there are a number of issues I have with the use of Make.

### Understanding

Magefiles are written in Go - a bonus for organisations and individuals that are already invested in Go.
Magefiles are just Go so can import any plain-old Go dependencies that you are familiar with.

Makefiles are esoteric (others may
disagree), but whilst the majority of us are comfortable writing and understanding basic Makefiles and Bash scripts, I've
found that the vast majority engineers (myself included) begin to struggle once things get a little more complex. 
that for myself and the vast majority engineers I've worked with begin to struggle when moving past basic examples.

### Don't repeat yourself

Magefiles allow importing of targets as vanilla Go imports, providing an easy way to share scripts throughout projets.

At larger organisations, the divergence of individual project Makefiles can turn into a maintenance nightmare. What starts
as a few projects with copied scripts can turn into a nightmare, with each project receives fixes and improvements in
isolation. Makefiles provide no solution here, and whilst there are many options to solve this problem with Make is to
add additional complexity. 

### Dependency management
     
Magefiles have a zero-install option, running with `go run`, requiring nothing but a Go install.
Magefiles can also be compiled to a binary and ran even without Go.

Often when writing Makefiles, anything past basic use-cases requires some kind of dependency, think `jq` to parse JSON
files, `docker` to run containers, `git` to manipulate tags as in this example. This means you need to also write bash
scripts to install these dependencies for each environment your scripts run in. Whilst not a show-stopper, it's nice to
remove these dependencies where possible and with Go being as popular as it is, there's a bunch of well-maintained
libraries (see `go-git`)
 
