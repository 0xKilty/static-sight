var color = false

var toggle_color = function () {
    setColor(color)
    color = !color
}

var pastColor = () => {
    var retrievedBoolean = localStorage.getItem("iankilty-com-color-mode");
    setColor(retrievedBoolean === "true")
}

var setColor = (bool) => {
    localStorage.setItem("iankilty-com-color-mode", bool)
    var root = document.documentElement

    var newColors = bool ? ["#fff", "#000", "#f3f3f3"] : ["#1e2021", "rgb(193, 193, 193)", "#282a36"]

    root.style.setProperty('--background-color', newColors[0]);
    root.style.setProperty('--text-color', newColors[1]);
    root.style.setProperty('--hljs-background', newColors[2]);
}