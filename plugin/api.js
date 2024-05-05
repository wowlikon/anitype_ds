const api = 'http://localhost:878';
var wt_tryed = false;
var wt_title = "?";

function SendData(data) {
  fetch(api+'/set', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error(error));
}

function get_wt_title(path){
  var req = new XMLHttpRequest(); 
  
  req.open('GET', `https://anitype.fun/anime/${path[1].split("?")[0]}`, false);   
  req.send(null);
  
  if(req.status == 200) {
    var parser = new DOMParser();
    var page = parser.parseFromString(req.responseText, "text/html");
    wt_title = page.body.getElementsByClassName("anime_central_body_title")[0].innerText;
  }

  wt_tryed = true;
}

function update() {
  var text = "";
  var count = "";
  var wt_url = "";
  var user_url = "";

  var url = window.location.href;
  var path = url.substring(20).split("/");
  switch (path[0]) {
    case "welcome":
      text = "На главной странице";
      break;

    case "anime":
      var title = document.body.getElementsByClassName("anime_central_body_title")[0];
      text = `Читает описание аниме "${title.innerText}"`;
      break;
      
    case "play":
      var title = document.body.getElementsByClassName("player_page_info_title")[0];
      text = `Смотрит аниме "${title.innerText}"`;
      break;
        
    case "watch_together":
      if (!wt_tryed) get_wt_title(path);
      text = `Совместный просмотр аниме "${wt_title}"`;
      count = document.body.getElementsByClassName("watch_together_hud_line")[0].childElementCount -1;
      wt_url = url;
      break;
    
    case "library":
      text = "Просматривает библиотеку";
      break;
      
    case "open":
      var list = path[1].split("?")[0];
      text = `Просматривает список "${decodeURI(list)}"`;
      break;
          
    case "settings":
      text = "Открыл настройки";
      break;

    default:
      if (path[0].startsWith("@")) {
        text = `Ссылка на профиль`;
        user_url = url.split("?")[0];
      }
      break;
    }
  if (text != "") {
    const setData = {
      wt: wt_url, usr: user_url,
      usrCount: count, text: text,
    };
    SendData(setData);
  }
}

window.onload = () => {setInterval(update, 5000)};