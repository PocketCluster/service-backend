# PocketCluster Index Deployment Log

- 10/10/2017
  * bolt -> coreos/bbolt for search stat/ supplement data/ social sharing
  * syntax search activated
  * ScalaJS based frontend removed (We'd need to move to npm)
  * detailed error log print

- 01/15/2017
  * Repository Meta/ Supplementary (Release + Tag) update internalized

- 01/13/2017
  * 5 Latest release note including tag is displayed

- 01/03/2017
  * Spinning indicator + page termination message added


## TODO

- [x] watch log and find `readme.html` failure with insufficient permission
- [ ] find 404 not found in readme.html from email -> fix
- [ ] show important author's profile in their title
- [ ] redirect page to a new, moved repository url
- [ ] prevent meta update fixes back the title
- [ ] deploy jQuery or a better alternative to dashboard app
- [ ] mark any update error to dashboard
- [ ] 404 error handling
- [ ] `/robots.txt` request doesn't go to nginx. why?
- [ ] sometimes, two identical repositories appear in the first page
- [ ] cannot click 'PocketCluster Index' to go back to home
- [ ] catch signal to stop as a service
- [ ] [LTG] rebuild the js app with `vue.js` + `infscr.js` + `mansonry` in mind that it should feel smooth
  * update infscr.js [`infinite-scroll`](https://infinite-scroll.com/) and masonry.
  * there is a glitch in paging index.html and categories as infscrl.js automatically increase page number
