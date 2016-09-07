# PocketCluster Index Frontend w/ Scala.js

### How to setup and live debug  

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