'use strict';

class SearchPage {
  constructor() {
    var grid = $('#indexes').masonry({
        columnWidth: '.card',
        itemSelector: '.card',
        percentPosition: true,
        transitionDuration: 10,
    });

    grid.imagesLoaded().progress(() => grid.masonry('layout'));

    grid.infinitescroll({
        loading: {
            selector: '#loading',
            finishedMsg: '<em>Feeling like adding your project? We have seats for you! Tweet <a href="https://twitter.com/stkim1">@stkim1</a> now!<em>',
            msg: $('<div id="infscr-loading" class="span12"><img alt="Loading..." src="' + staticRoot + '/img/default.svg"/></div>'),
        },
        navSelector: '#pagination',
        nextSelector: '#pagination a.next',
        itemSelector: '.card',
        path: function(currPage) {
            return $('#pagination a.next')[0].href + "&page=" + (currPage - 1);
        },
    },
    (elem) => {
        grid.masonry('appended', elem);
        grid.imagesLoaded().progress(() => grid.masonry('layout'));
    });

    // search form action
    $("form#nav-search").on('submit', e => this.searchSubmit(e));
  }

  // search request submit
  searchSubmit(evt) {
    evt.preventDefault();
    var term = $("form#nav-search input.form-control")[0].value.toString();
    if (term.length === 0) {
      return false;
    }
    var shref = "//" + window.location.host + "/search?term=" + encodeURI(term)
    window.location.replace(shref);
  }
}

document.addEventListener('DOMContentLoaded', () => new SearchPage());
