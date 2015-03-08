mokr
====

Basic tool to compile a [12factor](http://12factor.net/)/[heroku](https://devcenter.heroku.com/categories/language-support) app git repo into a docker container using [buildpacks](https://devcenter.heroku.com/articles/buildpacks).

Will interact with mesos/marathon eventually.

How
---
```
$ git clone https://github.com/deis/example-java-jetty.git
$ cd example-java-jetty
$ mokr myname/myimage
...
<snipped>
...
2015/03/08 20:15:46 image built: myname/myimage:a955d82
For a web app try:
    docker run -ti -p 8080:8080 -e PORT=8080 myname/myimage:a955d82 start web
```

Requires:
* Docker >= 1.5 to be installed on the users path
* Git >= 1.8.3 to be installed on the users path
* User able to run docker without sudo/admin rights

For other language example repos, see: https://github.com/deis?query=example

Huh?
----
1. grabs some details about the current commit (author, sha1, branch name)
2. performs a `git export` of the current commit and starts up a builder container
3. checks your code against a list of build packs to find which to compile with
4. vendors required binaries, downloads dependencies, compiles your code
5. creates a slug containing binaries and compiled code
6. creates a Dockerfile which adds in your slug to a prebuilt image
7. builds a docker container using the Dockerfile
