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

    listParam["auto"] = (listParam["auto"] != undefined ? listParam["auto"] : true);
    listParam["type"] = (listParam["type"] != undefined ? listParam["type"] : 'videos');
    console.log(listParam["id"])
    listParam["id"] = ((listParam["id"] != undefined) && (listParam["id"] != (NaN || "NaN")) ? parseInt(listParam["id"]) : rand(1, getCounts('videos')));
    return listParam
}

function rand(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function toggleAutoMode() {
    GlobaleState.autoMode = !GlobaleState.autoMode;
    setVideo()
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
    mediaFrame.onclick(randomImage())

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
            randomVideo();
            break;
        case "KeyQ":
        case "ArrowLeft":
            setMedia(GlobaleState.mediaID - 1);
            break;
        case "KeyD":
        case "ArrowRight":
            setMedia(GlobaleState.mediaID + 1);
            break;
        case "KeyA":
            toggleAutoMode();
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
    if (!(/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent)
        || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0, 4)))) {
        alert("Shortcuts\nRandom : R\nPrevious/Next : Left/Right\nAuto-Mode : A")
    }
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
