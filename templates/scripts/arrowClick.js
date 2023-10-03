var arrowClick = function (element) {
    var orientation = element.classList[1]
    if (orientation == "right") {
        element.classList.remove("right")
        element.classList.add("down")
        element.parentNode.parentNode.children[2].hidden = false
    } else {
        element.classList.remove("down")
        element.classList.add("right")
        element.parentNode.parentNode.children[2].hidden = true
    }
}