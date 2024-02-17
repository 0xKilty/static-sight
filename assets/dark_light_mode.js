var color = false

var toggle_color = function () {
    set_color(color)
    color = !color
}

var past_color = () => {
    set_color(localStorage.getItem("iankilty-com-color-mode") == "true")
}

var set_color = (bool) => {
    localStorage.setItem("iankilty-com-color-mode", bool)
    var root = document.documentElement

    const newColors = bool ? ["#fff", "#000", "#f3f3f3"] : ["#1e2021", "rgb(193, 193, 193)", "#282a36"];
    ['background-color', 'text-color', 'hljs-background'].forEach((prop, index) => root.style.setProperty(`--${prop}`, newColors[index]));

    var elements = document.getElementsByClassName("color_mode");

    for (var i = 0; i < elements.length; i++) {
        elements[i].innerHTML = bool ? "Dark" : "Light";
    }
}