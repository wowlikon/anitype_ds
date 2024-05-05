const api = 'http://localhost:878';
const useBtn = document.getElementById("useBtn");
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
    if (data.status) {
      useBtn.style = "background: var(--active);";
      console.log('Variable set to true');
      useBtn.innerText = "Включено";
      isActive = true;
    } else {
      useBtn.style = "background: var(--inactive);";
      console.log('Variable not set');
      useBtn.innerText = "Выключено";
      isActive = false;
    }
  })
  .catch((error) => {
    console.error('There has been a problem with your fetch operation:', error);
    useBtn.style = "background: var(--default);";
    useBtn.innerText = "Оффлайн";
    isActive = false;
  });
}

updBtn();
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