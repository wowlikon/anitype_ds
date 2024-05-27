const api = 'http://localhost:878';
var wt_try_genres = false;
var wt_try_tile = false;
var wt_title = "?";
var genres = "";

function SendData(url, data) {
  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => {console.error(url, data, error)});
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

  wt_try_tile = true;
}

function get_genres(path){
  var req = new XMLHttpRequest(); 

  req.open('GET', `https://anitype.fun/anime/${path[1].split("?")[0]}`, false);   
  req.send(null);
  
  if (req.status == 200) {
    var parser = new DOMParser();
    var page = parser.parseFromString(req.responseText, "text/html");
    keys = page.body.getElementsByClassName("anime_info_el_key");
    values = page.body.getElementsByClassName("anime_info_el_value");
    for (index = 0; index < keys.length; ++index) {
      if (keys[index].innerText == "Жанры") {
        genres = values[index].innerText;
        break;
      }
    }
  }
  wt_try_genres = true;
}

function update(url) {
  var count = 0;
  var text = "";
  var wt_url = "";
  var user_url = "";

  var path = url.substring(20).split("/");
  switch (path[0]) {
    case "welcome":
      text = "На главной странице";
      break;

    case "anime":
      var title = document.body.getElementsByClassName("anime_central_body_title")[0];
      genres = document.body.getElementsByClassName("anime_info_el_value")[4].innerText;
      text = `Читает описание аниме "${title.innerText}"`;
      break;
      
    case "play":
      var title = document.body.getElementsByClassName("player_page_info_title")[0];
      text = `Смотрит аниме "${title.innerText}"`;
      if (!wt_try_genres) get_genres(path);
      break;
        
    case "watch_together":
      if (!wt_try_tile) get_wt_title(path);
      if (!wt_try_genres) get_genres(path);
      text = `Совместный просмотр аниме "${wt_title}"`;
      count = document.body.getElementsByClassName("watch_together_hud_line")[0].childElementCount -1;
      wt_url = url;
      break;
    
    case "library":
      text = "Просматривает библиотеку";
      break;
    
    case "event":
      text = "Просматривает событие";
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
  if ((text != "") && window.location.href.startsWith("https://anitype.fun/")) {
  const setData = {
      genres: genres, wt: wt_url,
      usr: user_url, usrCount: count, text: text,
    };
    SendData(api+'/set', setData);
  }
  return text
}

if (window.location.href.startsWith("https://anitype.fun/")){
  window.onload = () => {setInterval(() => {update(window.location.href)}, 5000)};
}