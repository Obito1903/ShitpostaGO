var autoEnable = false
var GlobaleState = {
    mediaType: "video",
    mediaID: 1,
    autoMode: false
}

// var shitBase = window.location.origin
var shitBase = window.location.protocol + "//" + window.location.host
console.log(window.location.host)
function getParams() {
    var listParam = new Object();
    // get everything after "?" in the uri
    param = window.location.search.slice(1, window.location.search.length);
    // break down param key and value
    first = param.split("&");
    for (i = 0; i < first.length; i++) {
        second = first[i].split("=");
        listParam[second[0]] = second[1];
    }

    listParam["auto"] = (listParam["auto"] == !undefined ? listParam["auto"] : false);
    listParam["type"] = (listParam["type"] == !undefined ? listParam["type"] : 'videos');
    listParam["id"] = (parseInt(listParam["id"]) == NaN ? parseInt(listParam["id"]) : 1);
    return listParam
}

function rand(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function toggleAutoMode() {
    GlobaleState.autoMode = !GlobaleState.autoMode;
    updatePage()
}

function switchMediaType() {
    var button = document.getElementById('mediaType')
    if (GlobaleState.mediaType == "videos") {
        button.innerHTML == "Images"
        GlobaleState.mediaType = "images"
    } else {
        button.innerHTML == "Videos"
        GlobaleState.mediaType = "videos"
    }
    updatePage()
}

function getNewParams() {
    return "?type=" + GlobaleState.mediaType + "&id=" + GlobaleState.mediaID + "&auto=" + GlobaleState.autoMode
}

function updatePage() {
    var newURL = window.location.origin + window.location.pathname + getNewParams()
    var newTitle = toString(GlobaleState.mediaID)
    let newState = { additionalInformation: 'Changed Media' };

    // This will create a new entry in the browser's history, without reloading
    window.history.pushState(newState, newTitle, newURL);
    updateButtons();
}

function updateButtons() {
    var MediaTypeButton = document.getElementById('mediaType')
    var AutoButton = document.getElementById('autoButton')
    if (GlobaleState.autoMode) {
        AutoButton.classList.remove("bg-red-600", "hover:bg-red-700")
        AutoButton.classList.add("bg-green-600", "hover:bg-green-700")
    } else {
        AutoButton.classList.remove("bg-green-600", "hover:bg-green-700")
        AutoButton.classList.add("bg-red-600", "hover:bg-red-700")
    }
    if (GlobaleState.mediaType == "images") {
        MediaTypeButton.classList.remove("bg-green-600", "hover:bg-green-700")
        MediaTypeButton.classList.add("bg-red-600", "hover:bg-red-700")
        MediaTypeButton.innerHTML == "Images"
    } else {
        MediaTypeButton.classList.remove("bg-red-600", "hover:bg-red-700")
        MediaTypeButton.classList.add("bg-green-600", "hover:bg-green-700")
        MediaTypeButton.innerHTML == "Videos"
    }
}

function setVideo() {
    var mediaContainer = document.getElementById("media")
    mediaContainer.innerHTML = ""
    var mediaFrame = document.createElement("video");
    mediaFrame.setAttribute("autoplay", "");
    mediaFrame.setAttribute("controls", "");
    mediaFrame.setAttribute("name", GlobaleState.mediaType + GlobaleState.mediaID)
    mediaFrame.className = "h-screen w-screen"
    if (GlobaleState.autoMode) {
        mediaFrame.addEventListener("ended", randomVideo, false);
    }
    mediaFrame.id = "mediaFrame"
    mediaFrame.setAttribute("src", shitBase + "/shit?id=" + GlobaleState.mediaID + "&type=videos");
    mediaContainer.appendChild(mediaFrame)

    updatePage()
}

function setImage() {

    var mediaContainer = document.getElementById("media")
    mediaContainer.innerHTML = ""
    var mediaFrame = document.createElement("img");
    mediaFrame.setAttribute("name", GlobaleState.mediaType + GlobaleState.mediaID)
    mediaFrame.className = "max-w-full h-full"
    mediaFrame.id = "mediaFrame"
    mediaFrame.setAttribute("src", shitBase + "/shit?id=" + GlobaleState.mediaID + "&type=images");
    mediaContainer.appendChild(mediaFrame)

    updatePage()
}

function setMedia(ID) {
    console.log(GlobaleState)
    GlobaleState.mediaID = ID
    switch (GlobaleState.mediaType) {
        case "videos":
            setVideo(ID);
            break;
        case "images":
            setImage(ID);
        default:
            break;
    }
}

function processKey(key) {
    switch (key.code) {
        case "KeyS":
        case "KeyR":
            randomVideo()
            break;
        case "KeyQ":
        case "ArrowLeft":
            setMedia(GlobaleState.mediaID - 1)
            break;
        case "Keyd":
        case "ArrowRight":
            setMedia(GlobaleState.mediaID + 1)
            break;
        default:
            console.log(` ${key.code}`);
            break;
    }
}

function getCounts(type) {
    var count;

    xhttp = new XMLHttpRequest();

    xhttp.onreadystatechange = function () {
        if ((this.readyState == 4) && (this.status == 200)) {
            count = this.responseText
        }
    };

    xhttp.open("POST", shitBase + "/shitCount", false);
    xhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhttp.send("type=" + type);
    return count
}

function randomVideo() {
    var params = getParams();
    if ("min" in params) {
        var ID = rand(params["min"], getCounts(GlobaleState.mediaType)) + "&min=" + params["min"];
    } else {
        var ID = rand(1, getCounts(GlobaleState.mediaType));
    }
    setMedia(ID)
}

function randomImage() {
    window.location = "?type=images&auto=" + autoEnable + "&id=" + rand(1, nbImg);
}


function onload() {
    params = getParams();
    GlobaleState.autoMode = params["auto"]
    GlobaleState.mediaType = (params["type"] == !undefined ? params["auto"] : 'videos')
    GlobaleState.mediaID = parseInt(params["id"])
    console.log("yo")
    document.addEventListener('keydown', processKey)
    updateButtons();
    setMedia(GlobaleState.mediaID)
    document.addEventListener('swiped-right', function () {
        randomVideo()
    });
}
