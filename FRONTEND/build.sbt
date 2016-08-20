name := "INDEX-FRONT"

version := "1.0"

scalaVersion := "2.11.8"

libraryDependencies ++= Seq(
  "org.scala-js" %%% "scalajs-dom" % "0.9.0"
  ,"be.doeraene" %%% "scalajs-jquery" % "0.9.0"
)

//  This is to include all the JS dependencies in one file. We don't do this for now.
skip in packageJSDependencies := false
jsDependencies ++= Seq(
  "org.webjars" % "jquery" % "2.1.4" / "2.1.4/jquery.js"
  ,RuntimeDOM
)

// this is a hack to workaround http://stackoverflow.com/questions/31552605/intellij-sbt-sbt-native-packager-and-enableplugins-error
lazy val root = (project in file(".")).
  enablePlugins(ScalaJSPlugin).
  settings(
    name := "Scala.js Tutorial",
    scalaVersion := "2.11.8",
    version := "1.0"
  )