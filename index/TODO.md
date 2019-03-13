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


### SEARCH APP `JS` & `GO` PARTS
```js
$scope.processResults = function(data) {
    $scope.errorMessage = null;
    $scope.results = data;
    for(var i in $scope.results.hits) {
            hit = $scope.results.hits[i];
            hit.roundedScore = $scope.roundScore(hit.score);
            hit.explanationString = $scope.expl(hit.explanation);
            hit.explanationStringSafe = $sce.trustAsHtml(hit.explanationString);
            for(var ff in hit.fragments) {
                fragments = hit.fragments[ff];
                newFragments = [];
                for(var ffi in fragments) {
                    fragment = fragments[ffi];
                    safeFragment = $sce.trustAsHtml(fragment);
                    newFragments.push(safeFragment);
                }
                hit.fragments[ff] = newFragments;
            }
    }
    $scope.results.roundTook = $scope.roundTook(data.took);
};

$scope.searchTerm = function() {
    $http.post('/api/search', {
        "size": 10,
        "explain": true,
        "highlight":{},
        "query": {
            "term": $scope.term,
            "field": $scope.field,
        }
    }).
    success(function(data) {
        $scope.processResults(data);
    }).
    error(function(data, code) {

    });
};

$scope.roundTook = function(took) {
    if (took < 1000 * 1000) {
        return "less than 1ms";
    } else if (took < 1000 * 1000 * 1000) {
        return "" + Math.round(took / (1000*1000)) + "ms";
    } else {
        roundMs = Math.round(took / (1000*1000));
        return "" + roundMs/1000 + "s";
    }
};

updateFieldNames = function() {
    $http.get('/api/fields').success(function(data) {
        $scope.fieldNames = data.fields;
    }).
    error(function(data, code) {

    });
};
updateFieldNames();
```

```go
blserv.RegisterIndexName("repo", app.SearchIndex)
searchHandler := blserv.NewSearchHandler("repo")
listFieldsHandler := blserv.NewListFieldsHandler("repo")
router.Handle("/api/fields", listFieldsHandler).Methods("GET")
blserv "github.com/blevesearch/bleve/http"

type SearchRequest struct {
    Query            query.Query       `json:"query"`
    Size             int               `json:"size"`
    From             int               `json:"from"`
    Highlight        *HighlightRequest `json:"highlight"`
    Fields           []string          `json:"fields"`
    Facets           FacetsRequest     `json:"facets"`
    Explain          bool              `json:"explain"`
    Sort             search.SortOrder  `json:"sort"`
    IncludeLocations bool              `json:"includeLocations"`
}
```

### old infinite-scroll + masonry

```js
<script type="text/javascript">
    $(document).ready(function() {
        var grid = $('#indexes').masonry({
            columnWidth: '.card',
            itemSelector: '.card',
            percentPosition: true,
            transitionDuration: 10,
        });
        grid.imagesLoaded().progress( function() {
            grid.masonry('layout');
        });
        grid.infinitescroll({
            loading: {
                selector: '#loading',
                finishedMsg: '<em>Feeling like adding your project? We have seats for you! Tweet <a href="https://twitter.com/stkim1">@stkim1</a> now!<em>',
                msg: $('<div id="infscr-loading" class="span12"><img alt="Loading..." src="{{THEME_LINK}}/img/default.svg"/></div>'),
            },
            navSelector: '#pagination',
            nextSelector: '#pagination a.next',
            itemSelector: '.card',
        },
        function(elem) {
            grid.masonry('appended', elem);
            grid.imagesLoaded().progress( function() {
                grid.masonry('layout');
            });
        });
    });
</script>
```