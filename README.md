mokr
====

Basic tool to compile a [12factor](http://12factor.net/)/[heroku](https://devcenter.heroku.com/categories/language-support) app git repo using [buildpacks](https://devcenter.heroku.com/articles/buildpacks) into a docker container.

Will interact with mesos/marathon eventually.

```
$ git clone https://github.com/deis/example-java-jetty.git
$ cd example-java-jetty
$ mokr myname/myimage
...
<snipped>
...
2015/03/08 20:15:46 image built: myname/myimage:a955d82
2015/03/08 20:15:46 For a web app try:
2015/03/08 20:15:46     docker run -ti -p 8080:8080 -e PORT=8080 myname/myimage:a955d82 start web
```

Requires:
* Docker >= 1.5 to be installed on the users path
* Git >= 1.8.3 to be installed on the users path
* User able to run docker without sudo/admin rights

For other language example repos, see: https://github.com/deis?query=example
