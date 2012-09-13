mgomodel is a GO package that lets you interact with mongodb in a familiar model based way. It 
builds on top of the awsome mgo mongo driver http://labix.org/mgo.

mgomodel exposes the familiar save, delete, update, find, load methods one would expect with 
a normal ORM type setup. There is also a configurable mechanism for building in custom data validation.

Check out mgomodel_test.go for some examples. Docs are here: http://go.pkgdoc.org/github.com/mdennebaum/mgomodel