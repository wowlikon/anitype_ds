const useBtn = document.getElementById("useBtn");
var input = document.getElementById("genre");
var isActive = false;

function updBtn() {
  useBtn.style = "background: var(--default);";

  fetch(api+'/get', {
    method: 'GET',
    mode: 'cors',
  })
  .then((response) => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then((data) => {
    if (!data.hidden) {
      if (data.active) {
        useBtn.style = "background: var(--active);";
        useBtn.innerText = "Включено";
        isActive = true;
      } else {
        useBtn.style = "background: var(--inactive);";
        useBtn.innerText = "Выключено";
        isActive = false;
      }
    } else {
      useBtn.style = "background: var(--default);";
      useBtn.innerText = "Скрыто";
      isActive = false;
    }
    document.getElementById("label").outerHTML = '<p id="label"></p>'
  })
  .catch((error) => {
    console.error('There has been a problem with your fetch operation:', error);
    useBtn.style = "background: var(--default);";
    useBtn.innerText = "Оффлайн";
    document.getElementById("label").outerHTML = '<div id="label"><a href=\"https://github.com/wowlikon/anitype_ds/releases/tag/app\">Приложение не найдено :(</a></div>';
    isActive = false;
  });
}

function updBlock() {
  filter = document.getElementById("filters");
  filter.innerText = "";

  fetch(api+'/get_block', {
    method: 'GET',
    mode: 'cors',
  })
  .then((response) => {
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  })
  .then((data) => {
    arr = "";
    data.genres.map( function(item) {
      arr += `<li><label for="${item}">${item}</label></li>`;
    })
    filter.innerHTML = arr;
    document.querySelectorAll("label").forEach((item) => {
      item.addEventListener('click', () => {
        var blockData = {genre: item.innerText};
        SendData(api+'/del_block', blockData);
        updBlock();
      });
    });
  })
  .catch((error) => {
    console.error('There has been a problem with your fetch operation:', error);
  });
}

input.addEventListener("keypress", function(event) {
  if (event.key === "Enter") {
    event.preventDefault();
    if (input.value != "") {
      var blockData = {genre: input.value};
      SendData(api+'/add_block', blockData);
      input.value = "";
      updBlock();
    }
  }
});

updBtn();
updBlock();
useBtn.addEventListener("click",() => {
  if (isActive) {
    fetch(api+'/disenabled', {
      method: 'GET',
      mode: 'no-cors'
    });
  } else {
    fetch(api+'/enable', {
      method: 'GET',
      mode: 'no-cors'
    });
  }
  updBtn();
});