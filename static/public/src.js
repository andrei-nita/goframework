var section1 = {};
(function (o) {
    o.hello = function () {
        var hello = "Hello from section 1"
        console.log(hello)
    };
})(section1);
var section2 = {};
(function (o) {
    o.hello = function () {
        var hello = "Hello from section 2!!"
        console.log(hello)
    };
})(section2);
var ul1 = {};
(function (o) {
    o.hello = function () {
        var hello = "Hello from ul1!!"
        console.log(hello)
    };
})(ul1);
var footer = {};
(function (o) {
    o.hello = function () {
        var hello = "Hello from footer!!"
        console.log(hello)
    };
})(footer);
function toggleMenu() {
    var x = document.getElementById("myTopnav");
    if (x.className === "topnav") {
        x.className += " responsive";
    } else {
        x.className = "topnav";
    }
}

document.getElementById("menuBtn").addEventListener("click", toggleMenu);