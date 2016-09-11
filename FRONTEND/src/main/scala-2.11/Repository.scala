import org.scalajs.dom.raw.FormData
import org.scalajs.dom.{Element, Event, Node, ProgressEvent, XMLHttpRequest, html}
import org.scalajs.jquery.jQuery

import scala.collection.mutable
import scala.scalajs.js
import scala.scalajs.js.{Dictionary, JSApp, JSON}
import scala.scalajs.js.annotation.JSExport
import scala.scalajs.js.Dynamic.global
import scala.scalajs.js.JSConverters.JSRichGenMap

object Repository extends JSApp {
    def appendPar(targetNode: Node, text: String): Unit = {
        jQuery("body").append("<p>[" + text + "]</p>")
    }

    //@JSExport
    def addClickedMessage(): Unit = {
        jQuery("body").append("<p>You clicked the button!</p>")
    }

    def setupUI(): Unit = {
        jQuery("#click-me-button").click(addClickedMessage _)
        jQuery("body").append("<p>Hello World</p>")
    }

    @JSExport
    def previewRepository(click: Event) : Boolean = {

        jQuery("input#add-repo-desc").parent().removeClass("has-error")

        if (jQuery("input#add-repo-url").value().toString.length == 0) {
            jQuery("input#add-repo-url").parent().addClass("has-error")
            return false
        } else {
            jQuery("input#add-repo-url").parent().removeClass("has-error")
        }

        val requests = mutable.Map[String, String]()
        requests("add-repo-url")      = jQuery("input#add-repo-url").value().toString
        requests("add-repo-category") = jQuery("select#add-repo-category option:selected").value().toString
        requests("add-repo-desc")     = jQuery("input#add-repo-desc").value().toString
        requests("add-project-page")  = jQuery("input#add-project-page").value().toString
        requests("add-repo-logo")     = jQuery("input#add-repo-logo").value().toString

        val formdata = JSON.stringify(requests.toJSDictionary)

        val xhr = new XMLHttpRequest()
        xhr.open("POST", "/pocketcluster/dashboard/repository/preview")
        xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
        xhr.onload = { (e: Event) =>
            if (xhr.status == 200) {
                try {
                    val results = JSON.parse(xhr.responseText).asInstanceOf[Dictionary[String]].toMap[String, String]
                    println(results)
                    jQuery("input#add-repo-id").value(results("add-repo-id").toString)
                    jQuery("input#add-repo-title").value(results("add-repo-title").toString)
                    jQuery("input#add-repo-slug").value(results("add-repo-slug").toString)
                    jQuery("input#add-repo-desc").value(results("add-repo-desc").toString)
                } catch {
                    case e:Exception => println(e.printStackTrace())
                }
            }
        }
        xhr.onprogress = { (prog: ProgressEvent) => }
        xhr.ontimeout = { (e: Event) =>
            println("timeout!!!")
        }
        xhr.send(formdata)
        return false
    }

    @JSExport
    def submitRepository() : Boolean = {

        var critical_error = false
        if (jQuery("input#add-repo-url").value().toString.length == 0) {
            jQuery("input#add-repo-url").parent().addClass("has-error")
            critical_error = true
        } else {
            jQuery("input#add-repo-url").parent().removeClass("has-error")
        }
        if (jQuery("input#add-repo-desc").value().toString.length == 0) {
            jQuery("input#add-repo-desc").parent().addClass("has-error")
            critical_error = true
        } else {
            jQuery("input#add-repo-desc").parent().removeClass("has-error")
        }

        if (critical_error)
            return false

        val requests = mutable.Map[String, String]()
        requests("add-repo-url")      = jQuery("input#add-repo-url").value().toString
        requests("add-repo-category") = jQuery("select#add-repo-category option:selected").value().toString
        requests("add-repo-desc")     = jQuery("input#add-repo-desc").value().toString
        requests("add-project-page")  = jQuery("input#add-project-page").value().toString
        requests("add-repo-logo")     = jQuery("input#add-repo-logo").value().toString
        requests("add-repo-title")    = jQuery("input#add-repo-title").value().toString
        requests("add-repo-slug")     = jQuery("input#add-repo-slug").value().toString

        val formdata = JSON.stringify(requests.toJSDictionary)
        global.console.log(formdata)

        val xhr = new XMLHttpRequest()
        xhr.open("POST", "/pocketcluster/dashboard/repository/submit")
        xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
        xhr.onload = { (e: Event) =>
            if (xhr.status == 200) {
                try {
                    val results = JSON.parse(xhr.responseText).asInstanceOf[Dictionary[String]].toMap[String, String]
                    jQuery("input#add-repo-id").value(results("add-repo-id").toString)
                    jQuery("input#add-repo-title").value(results("add-repo-title").toString)
                    jQuery("input#add-repo-slug").value(results("add-repo-slug").toString)
                    jQuery("input#add-repo-desc").value(results("add-repo-desc").toString)
                } catch {
                    case e:Exception => println(e.printStackTrace())
                }
            }
        }
        xhr.onprogress = { (prog: ProgressEvent) => }
        xhr.ontimeout = { (e: Event) =>
            println("timeout!!!")
        }

        xhr.send(formdata)
        return false
    }

    def main(): Unit = {
        //jQuery(setupUI _)
    }

}