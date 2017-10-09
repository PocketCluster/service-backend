'use strict';

// A Sticky Notes app.
class RepositoryManagerApp {

  // Initializes the Sticky Notes app.
  constructor() {

    // Shortcuts to DOM Elements.
    this.addRepoURL  = document.getElementById("add-repo-url");
    this.addRepoDesc = document.getElementById("add-repo-desc");
    this.addRepoCat  = document.getElementById("add-repo-cat");
    this.addRepoProj = document.getElementById("add-repo-proj");
    this.addRepoLogo = document.getElementById("add-repo-logo");

    this.isRepoDuplicated = false

    // Preview click event
    document.getElementById("btn-preview").addEventListener('click', () => this.preview());

    // Submit action event on form
    document.getElementById("add-repo-form").addEventListener('submit', e => this.submit(e));
  }

  // preview
  preview() {
    // reset errors
    this.addRepoDesc.parentNode.classList.remove("has-error")
    if (this.addRepoURL.value.toString().length == 0) {
      this.addRepoURL.parentNode.classList.add("has-error")
      return false
    } else {
      this.addRepoURL.parentNode.classList.remove("has-error")
    }

    let app    = this,
        rdesc  = document.getElementById("add-repo-desc"),
        rcat   = document.getElementById("add-repo-cat"),
        rproj  = document.getElementById("add-repo-proj"),
        rlogo  = document.getElementById("add-repo-logo"),
        repoid = document.getElementById("add-repo-id"),
        rtitle = document.getElementById("add-repo-title"),
        rslug  = document.getElementById("add-repo-slug"),
        sstat  = document.getElementById("add-repo-status"),
        srson  = document.getElementById("add-repo-status-reason");

    // create request
    let requests = {
      "add-repo-url":  this.addRepoURL.value.toString(),
      "add-repo-cat":  this.addRepoCat.options[this.addRepoCat.selectedIndex].value.toString(),
      "add-repo-proj": this.addRepoProj.value.toString(),
      "add-repo-logo": this.addRepoLogo.value.toString(),
    };

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/pocketcluster/dashboard/repository/preview")
    xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    xhr.onload = function () {
      if (xhr.status == 200) {
          var ret = JSON.parse(xhr.responseText)
          var results = Object.keys(ret).reduce((map, key) => map.set(key, ret[key]), new Map());
          // https://stackoverflow.com/questions/36644438/how-to-convert-a-plain-object-into-an-es6-map
          // const buildMap = o => Object.keys(o).reduce((m, k) => m.set(k, o[k]), new Map());

          repoid.value = results.get("add-repo-id")
          rtitle.value = results.get("add-repo-title")
          rslug.value  = results.get("add-repo-slug")
          rdesc.value  = results.get("add-repo-desc")

          switch(ret.status) {
            case "processed": {
              sstat.classList.remove("panel-warning")
              sstat.classList.add("panel-success")
              sstat.style.visibility = "visible"
              srson.textContent = "Everything went ok! Reload page!"
              break
            }
            case "duplicated": {
              app.isRepoDuplicated = true
              sstat.classList.remove("panel-success")
              sstat.classList.add("panel-warning")
              sstat.style.visibility = "visible"
              srson.textContent = ret.reason.toString()

              // update additional fields
              rcat.value  = results.get("add-repo-cat").toString()
              rproj.value = results.get("add-repo-proj").toString()
              rlogo.value = results.get("add-repo-logo").toString()
              break
            }
            case "error": {
              sstat.classList.remove("panel-success")
              sstat.classList.add("panel-danger")
              sstat.style.visibility = "visible"
              srson.textContent = ret.reason.toString()
              break
            }
          }
      }
    }
    xhr.onprogress = function (prog) {
      console.log(prog)
    }
    xhr.ontimeout = function() {
      console.log("timeout!!!")
    }
    xhr.send(JSON.stringify(requests))
    return false
  }

  // submit
  submit(evt) {
    evt.preventDefault();

    var critical_error = false
    if (this.addRepoURL.value.toString().length == 0) {
      this.addRepoURL.parentNode.classList.add("has-error")
      critical_error = true
    } else {
      this.addRepoURL.parentNode.classList.remove("has-error")
    }
    if (this.addRepoDesc.value.toString().length == 0) {
      this.addRepoDesc.parentNode.classList.add("has-error")
      critical_error = true
    } else {
      this.addRepoDesc.parentNode.classList.remove("has-error")
    }
    if (critical_error)
      return false

    let app   = this,
        rcat  = document.getElementById("add-repo-cat"),
        rproj = document.getElementById("add-repo-proj"),
        rlogo = document.getElementById("add-repo-logo"),
        sstat = document.getElementById("add-repo-status"),
        srson = document.getElementById("add-repo-status-reason");

    // create request
    let requests = {
        "add-repo-url":   document.getElementById("add-repo-url").value.toString(),
        "add-repo-cat":   rcat.options[rcat.selectedIndex].value.toString(),
        "add-repo-desc":  document.getElementById("add-repo-desc").value.toString(),
        "add-repo-proj":  document.getElementById("add-repo-proj").value.toString(),
        "add-repo-logo":  document.getElementById("add-repo-logo").value.toString(),
        "add-repo-title": document.getElementById("add-repo-title").value.toString(),
        "add-repo-slug":  document.getElementById("add-repo-slug").value.toString(),
    }

    var xhr = new XMLHttpRequest();
    if (this.isRepoDuplicated) {
        xhr.open("POST", "/pocketcluster/dashboard/repository/update")
    } else {
        xhr.open("POST", "/pocketcluster/dashboard/repository/submit")
    }
    xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    xhr.onload = function(evt) {
      if (xhr.status == 200) {
        var ret = JSON.parse(xhr.responseText)
        switch(ret.status) {
          case "processed": {
            app.isRepoDuplicated = true
            sstat.classList.remove("panel-warning")
            sstat.classList.add("panel-success")
            sstat.style.visibility = "visible"
            srson.textContent = "Everything went ok! Reload page!"
            break
          }
          case "duplicated": {
            app.isRepoDuplicated = true
            sstat.classList.remove("panel-success")
            sstat.classList.add("panel-warning")
            sstat.style.visibility = "visible"
            srson.textContent = ret.reason.toString()

            // update additional fields
            rcat.value  = results.get("add-repo-cat").toString()
            rproj.value = results.get("add-repo-proj").toString()
            rlogo.value = results.get("add-repo-logo").toString()
            break
          }
          case "error": {
            sstat.classList.remove("panel-success")
            sstat.classList.add("panel-danger")
            sstat.style.visibility = "visible"
            srson.textContent = ret.reason.toString()
            break
          }
        }
      }
    }
    xhr.onprogress = function (prog) {
      console.log(prog)
    }
    xhr.ontimeout = function() {
      console.log("timeout!!!")
    }
    xhr.send(JSON.stringify(requests))
    return false
  }
}

// On load start the app.
window.addEventListener('load', () => new RepositoryManagerApp());
