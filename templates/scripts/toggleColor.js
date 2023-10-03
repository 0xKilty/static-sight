var color = true;
var toggle_color = function () {
    const main = document.getElementById("main")
    if (color) {
        main.style.background = "#fff"
        main.style.color = "#000"
    } else {
        main.style.background = "#1e2021"
        main.style.color = "rgb(193, 193, 193)"
    }
    color = !color
}