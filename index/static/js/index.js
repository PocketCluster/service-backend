'use strict';

class IndexFrontPage {
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
            msg: $('<div id="infscr-loading" class="span12"><img alt="Loading..." src="/img/default.svg"/></div>'),
        },
        navSelector: '#pagination',
        nextSelector: '#pagination a.next',
        itemSelector: '.card',
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

  }

}

document.addEventListener('DOMContentLoaded', () => new IndexFrontPage());
