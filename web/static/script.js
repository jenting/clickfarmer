window.addEventListener('load', function () {
    updateClickLabels();
})

function colorClicked(color) {
    console.log("color clicked", color);
    // hmm
}

function updateClickLabels() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var response = JSON.parse(this.responseText);
            document.getElementById("color-label-red").innerHTML = "red: " + response.redClicks
            document.getElementById("color-label-green").innerHTML = "red: " + response.greenClicks
            document.getElementById("color-label-blue").innerHTML = "red: " + response.blueClicks
        }
    };
    xhttp.open("GET", "/api/clicks", true);
    xhttp.send();
}
