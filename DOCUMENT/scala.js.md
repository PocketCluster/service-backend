# PocketCluster Index Frontend w/ Scala.js

**Various Tips (09/12/2016)** 

1. `@JSExport` could be used to replace inline JS function  
2. `scala.scalajs.js.JSConverters.JSRichGenMap` package to required to convert `JSDictionary` <-> `Map`  
For more about Interoperability between `JS` and `Scala` take a look at [Interoperability](https://www.scala-js.org/doc/interoperability/) and [Scala.js](https://github.com/scala-js/scala-js)
3. Using `Future`, XHR connection could be more succinct.
4. `JSDictionary` to `Map` has some issues. Use following snippet to catch exceptions.

  ```scala
  try {
  } catch {
      case e:Exception => e.printStackTrace
  }
  ```
5. Interop between infinite scroll and Mansory might be possible in the future?

**Bootstrap + Mansory + Infinite Scroll + ImageLoaded (09/10/2016)**

Components  

  - [Bootstrap](http://getbootstrap.com/)
  - [Mansory](http://masonry.desandro.com/)
  - [JQuery Infinite Scroll](https://github.com/infinite-scroll/infinite-scroll)
  - [JQuery Image Loaded](http://imagesloaded.desandro.com/)

Tutorials

  - [Infinite Scrolling and Grids](https://www.sitepoint.com/jquery-infinite-scrolling-demos/)
  - [Getting Bootstrap Tabs to Play Nice with Masonry](https://www.sitepoint.com/bootstrap-tabs-play-nice-with-masonry/)
  - [Infinite Scroll + Masonry + ImagesLoaded](https://gist.github.com/gregrickaby/10383879)

Example
  - [Infinite Scroll Example](Infinite-scroll-example.zip)

**How to setup and live debug (09/07/2016)** 

1. Once you setup Scala.js with sbt, 

    ```sh
    $ sbt
    > run
    [info] Compiling 1 Scala source to (...)/scala-js-tutorial/target/scala-2.11/classes...
    [info] Running tutorial.webapp.TutorialApp
    Hello world!
    [success] (...)
    ```

2. Generate JavaScript

    To generate a single JavaScript file using sbt, just use the fastOptJS task:

    ```sh
    > fastOptJS
    [info] Fast optimizing (...)/scala-js-tutorial/target/scala-2.11/scala-js-tutorial-fastopt.js
    [success] (...)
    This will perform some fast optimizations and generate the target/scala-2.11/scala-js-tutorial-fastopt.js file containing the JavaScript code.
    ```

    Re-typing fastOptJS each time you change your source file is cumbersome. Luckily sbt is able to watch your files and recompile as needed:
    
    ```sh
    > ~fastOptJS
    [success] (...)
    1. Waiting for source changes... (press enter to interrupt)
    From this point in the tutorial we assume you have an sbt with this command running, so we don’t need to bother with rebuilding each time.
    ```

> Reference  
  
- <https://www.scala-js.org/tutorial/basic/>