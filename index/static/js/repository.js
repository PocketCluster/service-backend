'use strict';

class RepositoryPage {
  constructor() {
    // initialize social share
    SocialShareKit.init({url:window.location.href});

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
    window.location.assign(shref);
  }
}

document.addEventListener('DOMContentLoaded', () => new RepositoryPage());
